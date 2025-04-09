package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionDirectoryName struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionDirectoryName) Validate() error {
	if r.Value != "" {
		if r.Value != "lowercase-hyphens" {
			return fmt.Errorf("%s supports 'lowercase-hyphens' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleActionDirectoryName) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	if r.Value == "lowercase-hyphens" {
		regex := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
		m := regex.MatchString(a.DirName)
		if !m {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action directory name '%s' must be lower-case and hyphens only", a.DirName))
			return false, nil
		}
	}

	return true, nil
}

func (r RuleActionDirectoryName) GetConfigName() string {
	return r.ConfigName
}
