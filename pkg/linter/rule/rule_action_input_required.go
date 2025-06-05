package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// RuleActionInputRequired checks whether specific input attributes are defined (e.g. 'description').
// Currently, only the 'description' attribute is supported.
type RuleActionInputRequired struct {
	Value      []string
	ConfigName string
	IsError    []bool
}

func (r RuleActionInputRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "description" {
				return fmt.Errorf("%s can only contain 'description'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionInputRequired) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	for inputName, input := range a.Inputs {
		for i, v := range r.Value {
			if v == "description" && input.Description == "" {
				printErrOrWarn(r.ConfigName, r.IsError[i], fmt.Sprintf("action '%s' input '%s' does not have a required %s", a.DirName, inputName, v), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleActionInputRequired) GetConfigName() string {
	return r.ConfigName
}
