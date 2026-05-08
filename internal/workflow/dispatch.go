package workflow

// Dispatch represents a 'workflow_dispatch' field in a GitHub Actions workflow parsed from YAML.
type Dispatch struct {
	Inputs map[string]*Input `yaml:"inputs"`
}
