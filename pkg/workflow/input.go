package workflow

// Input represents an input of a GitHub Actions workflow parsed from YAML.
type Input struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}
