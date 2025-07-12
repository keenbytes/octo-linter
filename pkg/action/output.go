package action

// Output represents an output of a GitHub Actions action parsed from YAML.
type Output struct {
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}
