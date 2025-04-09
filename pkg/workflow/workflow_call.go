package workflow

type WorkflowCall struct {
	Inputs map[string]*WorkflowInput `yaml:"inputs"`
}
