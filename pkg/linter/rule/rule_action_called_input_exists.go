package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionCalledInputExists struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionCalledInputExists) Validate() error {
	return nil
}

func (r RuleActionCalledInputExists) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if r.Value {
		re := regexp.MustCompile(`\${{[ ]*inputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if a.Inputs == nil || a.Inputs[string(f[1])] == nil {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls an input '%s' that does not exist", a.DirName, string(f[1])))
				compliant = false
			}
		}
	}

	return
}

func (r RuleActionCalledInputExists) GetConfigName() string {
	return r.ConfigName
}
