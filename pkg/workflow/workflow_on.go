package workflow

type WorkflowOn struct {
	WorkflowCall     *WorkflowCall     `yaml:"workflow_call"`
	WorkflowDispatch *WorkflowDispatch `yaml:"workflow_dispatch"`
}
