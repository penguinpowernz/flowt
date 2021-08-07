package flowt

import "errors"

type Flow struct {
	Name    string
	Comment string
	StartAt string
	Steps   map[string]Step
}

func (flw Flow) Execute(ex *Execution) error {
	name := flw.StartAt
	for {
		step, found := flw.Steps[name]
		if !found {
			return errors.New("step not found: " + name)
		}

		stx, err := ExecuteStep(step, name, ex.State())
		ex.CompleteStep(stx)
		// TODO try to enter catch/retry step
		if stx.next == "$" {
			return err
		}

		name = stx.next
	}
}
