package workflow

import (
	"gopkg.pl/mikogs/octo-linter/pkg/step"
)

type WorkflowJob struct {
	Name   string            `yaml:"name"`
	Uses   string            `yaml:"uses"`
	RunsOn interface{}       `yaml:"runs-on"`
	Steps  []*step.Step      `yaml:"steps"`
	Env    map[string]string `yaml:"env"`
	Needs  interface{}       `yaml:"needs,omitempty"`
}

func (wj *WorkflowJob) SetParentType(t string) {
	for _, s := range wj.Steps {
		s.ParentType = t
	}
}
