package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

// RuleActionDirectoryName checks whether the action directory name adheres to the selected naming convention.
// Currently, only 'lowercase-hyphens' is supported, meaning the name must be entirely lowercase and use hyphens only.
type RuleActionDirectoryName struct {
	Value      string
	ConfigName string
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

func (r RuleActionDirectoryName) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if r.Value == "lowercase-hyphens" {
		regex := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
		m := regex.MatchString(a.DirName)
		if !m {
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("action directory name '%s' must be lower-case and hyphens only", a.DirName), chWarnings, chErrors)
			return false, nil
		}
	}

	return
}

func (r RuleActionDirectoryName) GetConfigName() string {
	return r.ConfigName
}
