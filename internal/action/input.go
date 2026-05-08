package action

// Input represents an input in a GitHub Actions action parsed from YAML.
type Input struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}
