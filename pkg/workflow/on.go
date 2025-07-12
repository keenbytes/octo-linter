package workflow

// On represents a 'on' field in a GitHub Actions workflow parsed from YAML.
type On struct {
	WorkflowCall     *Call     `yaml:"workflow_call"`
	WorkflowDispatch *Dispatch `yaml:"workflow_dispatch"`
}
