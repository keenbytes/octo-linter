package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionInputValue struct {
	Value      map[string]string
	ConfigName string
	LogLevel   int
	IsError    map[string]bool
}

func (r RuleActionInputValue) Validate() error {
	if len(r.Value) > 0 {
		for k, v := range r.Value {
			if k != "name" {
				return fmt.Errorf("%s can only contain 'name' key", r.ConfigName)
			}
			if v != "lowercase-hyphens" {
				return fmt.Errorf("%s supports 'lowercase-hyphens' or empty value only", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionInputValue) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	for inputName := range a.Inputs {
		for k, v := range r.Value {
			if k == "name" && v != "" {
				regex := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
				m := regex.MatchString(inputName)
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError[k], r.LogLevel, fmt.Sprintf("action '%s' input '%s' must be lower-case and hyphens only", a.DirName, inputName))
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func (r RuleActionInputValue) GetConfigName() string {
	return r.ConfigName
}
