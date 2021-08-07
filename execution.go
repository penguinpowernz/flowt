package flowt

import (
	"io"
	"io/ioutil"

	"github.com/pborman/uuid"
)

func NewExecutionStream(r io.Reader) (*Execution, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewExecution(data)
}

func NewExecution(input []byte) (*Execution, error) {
	exc := &Execution{
		ID:    uuid.New(),
		input: string(input),
	}

	return exc, nil
}

type Execution struct {
	ID           string `json:"id"`
	onStepChange func(string)
	input        string
	steps        []StepExecution
}

func (exc *Execution) CompleteStep(stx *StepExecution) {
	_stx := *stx
	_stx._state = nil
	exc.steps = append(exc.steps, _stx)
}

func (exc *Execution) State() string {
	if len(exc.steps) == 0 {
		return exc.input
	}

	return exc.steps[len(exc.steps)-1].output
}
