package workflow

// Call represents a 'workflow_call' field in a GitHub Actions workflow parsed from YAML.
type Call struct {
	Inputs map[string]*Input `yaml:"inputs"`
}
