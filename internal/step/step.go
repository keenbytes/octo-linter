// Package step contains code related to steps in GitHub Actions workflows and actions.
package step

// Step represents a GitHub Actions step parsed from a workflow or action file.
type Step struct {
	ParentType string
	Name       string            `yaml:"name"`
	ID         string            `yaml:"id"`
	Uses       string            `yaml:"uses"`
	Shell      string            `yaml:"bash"`
	Env        map[string]string `yaml:"env"`
	Run        string            `yaml:"run"`
	With       map[string]string `yaml:"with"`
}
