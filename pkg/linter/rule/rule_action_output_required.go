package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionOutputRequired struct {
	Value      []string
	ConfigName string
	LogLevel   int
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

func (r RuleActionOutputRequired) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	for outputName, output := range a.Outputs {
		for i, v := range r.Value {
			if v == "description" && output.Description == "" {
				printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("action '%s' output '%s' does not have a required %s", a.DirName, outputName, v))
			}
		}
	}

	return true, nil
}

func (r RuleActionOutputRequired) GetConfigName() string {
	return r.ConfigName
}
