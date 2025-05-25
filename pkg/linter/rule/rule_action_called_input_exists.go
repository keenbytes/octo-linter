package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

// RuleActionCalledInputExists scans the action code for all input references and verifies that each has been previously defined.
// During action execution, if a reference to an undefined input is found, it is replaced with an empty string.
type RuleActionCalledInputExists struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleActionCalledInputExists) Validate() error {
	return nil
}

func (r RuleActionCalledInputExists) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if r.Value {
		re := regexp.MustCompile(`\${{[ ]*inputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if a.Inputs == nil || a.Inputs[string(f[1])] == nil {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("action '%s' calls an input '%s' that does not exist", a.DirName, string(f[1])), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleActionCalledInputExists) GetConfigName() string {
	return r.ConfigName
}
