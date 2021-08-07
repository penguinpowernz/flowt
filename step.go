package flowt

// Step represents a step in the state machine
type Step struct {
	Type     string
	Resource string
	Next     string
	End      bool

	InputPath  string
	OutputPath string
	ResultPath string
	// Parameters interface{}
	// ResultSelector interface{}

	// pass state
	Result interface{}

	// choice state
	Choices []Choice
	Default string

	// error state
	Error string
	Cause string

	// wait state
	Seconds       int
	Timestamp     string
	TimestampPath string
	SecondsPath   string
}

// Execute will run the specified step with given execution
// func (step Step) Execute(ex *Execution) (string, error) {
// 	ex.buildStateInput(step)

// 	var next string = "$"
// 	var exe stateExecutor
// 	switch step.Type {
// 	case "Task":
// 		exe = taskExecutor(step)
// 	case "Fail":
// 		return next, errors.New(step.Error + ": " + step.Cause)
// 	case "Choice":
// 		exe = choiceExecutor(step)
// 	case "Succeed":
// 		return next, nil
// 	case "Wait":
// 		exe = waitExecutor(step)
// 	case "Pass":
// 		exe = passExecutor(step)
// 	}

// 	next, data, err := exe(x.)
// 	if err != nil {
// 		return "$", nil
// 	}

// 	err = ex.applyResultData(step, data)

// 	return next, err
// }
