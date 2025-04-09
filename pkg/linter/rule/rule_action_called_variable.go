package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionCalledVariable struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionCalledVariable) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleActionCalledVariable) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if r.Value == "uppercase-underscores" {
		reName := regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

		varTypes := []string{"env", "var", "secret"}
		for _, v := range varTypes {
			re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
			found := re.FindAllSubmatch(a.Raw, -1)
			for _, f := range found {
				m := reName.MatchString(string(f[1]))
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a variable '%s' that must be alphanumeric uppercase and underscore only", a.DirName, string(f[1])))
					compliant = false
				}
			}
		}
	}

	return
}

func (r RuleActionCalledVariable) GetConfigName() string {
	return r.ConfigName
}
