package step

type Step struct {
	ParentType string
	Name       string            `yaml:"name"`
	Id         string            `yaml:"id"`
	Uses       string            `yaml:"uses"`
	Shell      string            `yaml:"bash"`
	Env        map[string]string `yaml:"env"`
	Run        string            `yaml:"run"`
	With       map[string]string `yaml:"with"`
}
