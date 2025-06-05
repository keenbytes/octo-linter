package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// RuleActionRequired checks whether the specified action fields are present, e.g. 'name'.
type RuleActionRequired struct {
	Value      []string
	ConfigName string
	IsError    []bool
}

func (r RuleActionRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "name" && v != "description" {
				return fmt.Errorf("%s can only contain values of 'name' and/or 'details'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionRequired) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	for i, v := range r.Value {
		if (v == "name" && a.Name == "") || (v == "description" && a.Description == "") {
			printErrOrWarn(r.ConfigName, r.IsError[i], fmt.Sprintf("action '%s' does not have a required %s", a.DirName, v), chWarnings, chErrors)
			compliant = false
		}
	}

	return
}

func (r RuleActionRequired) GetConfigName() string {
	return r.ConfigName
}
