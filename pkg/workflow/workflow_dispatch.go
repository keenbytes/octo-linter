package workflow

type WorkflowDispatch struct {
	Inputs map[string]*WorkflowInput `yaml:"inputs"`
}
