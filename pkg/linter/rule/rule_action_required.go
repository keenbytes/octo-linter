package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionRequired struct {
	Value      []string
	ConfigName string
	LogLevel   int
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

func (r RuleActionRequired) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	for i, v := range r.Value {
		if (v == "name" && a.Name == "") || (v == "description" && a.Description == "") {
			printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("action '%s' does not have a required %s", a.DirName, v))
		}
	}

	return true, nil
}

func (r RuleActionRequired) GetConfigName() string {
	return r.ConfigName
}
