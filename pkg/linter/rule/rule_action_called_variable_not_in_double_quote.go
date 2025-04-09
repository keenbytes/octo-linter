package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionCalledVariableNotInDoubleQuote struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionCalledVariableNotInDoubleQuote) Validate() error {
	return nil
}

func (r RuleActionCalledVariableNotInDoubleQuote) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if r.Value {
		re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a variable '%s' that is in double quotes", a.DirName, string(f[1])))
			compliant = false
		}
	}

	return
}

func (r RuleActionCalledVariableNotInDoubleQuote) GetConfigName() string {
	return r.ConfigName
}
