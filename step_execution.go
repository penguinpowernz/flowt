package flowt

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
)

func ExecuteStep(step Step, name, input string) (*StepExecution, error) {
	stx := &StepExecution{name: name, input: input, start: time.Now()}
	defer func() { stx.end = time.Now() }()

	if err := stx.prepareState(step); err != nil {
		return stx, err
	}
	if err := stx.execute(step); err != nil {
		return stx, err
	}
	err := stx.prepareOutput(step)
	return stx, err
}

type StepExecution struct {
	name string

	input  string // input is the original input, before being transformed into the state by the step
	state  string // state is what is used by the step to do work
	output string // output is the output of the step using the state to do work
	next   string
	result []byte

	logs   io.ReadWriteCloser
	_state *gabs.Container

	start time.Time
	end   time.Time
	err   error
}

func (ex *StepExecution) execute(step Step) (err error) {
	quit := "$"
	var exe stateExecutor
	switch step.Type {
	case "Task":
		exe = taskExecutor(step)
	case "Fail":
		ex.next = quit
		return errors.New(step.Error + ": " + step.Cause)
	case "Choice":
		exe = choiceExecutor(step)
	case "Succeed":
		ex.next = quit
		return nil
	case "Wait":
		exe = waitExecutor(step)
	case "Pass":
		exe = passExecutor(step)
	}

	next, data, err := exe(ex._state)
	if err != nil {
		ex.next = quit
		return err
	}

	ex.next = next
	ex.result = data

	return
}

func (ex *StepExecution) prepareState(step Step) (err error) {
	ex._state, err = gabs.ParseJSON([]byte(ex.input))
	if err != nil {
		return
	}

	if step.InputPath != "" {
		inpPath := strings.TrimPrefix(step.InputPath, "$.")
		ex._state = ex._state.Path(inpPath)
	}

	// TODO: handle step.Parameters

	ex.state = ex._state.StringIndent("", "  ")
	return nil
}

func (ex *StepExecution) prepareOutput(step Step) error {
	if len(ex.result) == 0 {
		return nil
	}

	c, err := gabs.ParseJSON(ex.result)
	if err != nil {
		return err
	}

	// TODO: handle ResultsSelector

	if step.ResultPath == "" {
		ex._state.Merge(c)
	} else {
		resPath := strings.TrimPrefix(step.ResultPath, "$.")
		ex._state.Merge(c.Path(resPath))
	}

	if step.OutputPath != "" {
		outPath := strings.TrimPrefix(step.OutputPath, "$.")
		ex._state = ex._state.Path(outPath)
	}

	ex.output = ex._state.String()

	return nil
}
