package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// RuleActionOutputRequired checks whether specific output attributes are defined (e.g. 'description').
// Currently, only the 'description' attribute is supported.
type RuleActionOutputRequired struct {
	Value      []string
	ConfigName string
	IsError    []bool
}

func (r RuleActionOutputRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "description" {
				return fmt.Errorf("%s can only contain 'description'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionOutputRequired) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	for outputName, output := range a.Outputs {
		for i, v := range r.Value {
			if v == "description" && output.Description == "" {
				printErrOrWarn(r.ConfigName, r.IsError[i], fmt.Sprintf("action '%s' output '%s' does not have a required %s", a.DirName, outputName, v), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleActionOutputRequired) GetConfigName() string {
	return r.ConfigName
}
