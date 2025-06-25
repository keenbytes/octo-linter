package action

import (
	"github.com/keenbytes/octo-linter/v2/pkg/step"
)

type ActionRuns struct {
	Using string       `yaml:"using"`
	Steps []*step.Step `yaml:"steps"`
}

func (ar *ActionRuns) SetParentType(t string) {
	for _, s := range ar.Steps {
		s.ParentType = t
	}
}

func (ar *ActionRuns) GetStep(id string) *step.Step {
	for _, s := range ar.Steps {
		if s.Id == id {
			return s
		}
	}
	return nil
}
