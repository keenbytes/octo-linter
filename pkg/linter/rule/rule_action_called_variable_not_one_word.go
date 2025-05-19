package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionCalledVariableNotOneWord struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleActionCalledVariableNotOneWord) Validate() error {
	return nil
}

func (r RuleActionCalledVariableNotOneWord) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if r.Value {
		re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("action '%s' calls a variable '%s' that is invalid", a.DirName, string(f[1])), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleActionCalledVariableNotOneWord) GetConfigName() string {
	return r.ConfigName
}
