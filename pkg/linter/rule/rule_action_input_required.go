package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionInputRequired struct {
	Value      []string
	ConfigName string
	LogLevel   int
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

func (r RuleActionInputRequired) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	for inputName, input := range a.Inputs {
		for i, v := range r.Value {
			if v == "description" && input.Description == "" {
				printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("action '%s' input '%s' does not have a required %s", a.DirName, inputName, v))
			}
		}
	}

	return true, nil
}

func (r RuleActionInputRequired) GetConfigName() string {
	return r.ConfigName
}
