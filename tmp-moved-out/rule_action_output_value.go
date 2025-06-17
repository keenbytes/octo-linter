package rule

import (
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// RuleActionOutputValue verifies whether the action output field follows the specified naming convention —
// for example, ensuring the 'name' field uses 'lowercase-hyphens' (lowercase letters, digits, and hyphens only).
type RuleActionOutputValue struct {
	Value      map[string]string
	ConfigName string
	IsError    map[string]bool
}

func (r RuleActionOutputValue) Validate() error {
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

func (r RuleActionOutputValue) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	for outputName := range a.Outputs {
		for k, v := range r.Value {
			if k == "name" && v != "" {
				regex := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
				m := regex.MatchString(outputName)
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError[k], fmt.Sprintf("action '%s' output '%s' must be lower-case and hyphens only", a.DirName, outputName), chWarnings, chErrors)
					return false, nil
				}
			}
		}
	}

	return
}

func (r RuleActionOutputValue) GetConfigName() string {
	return r.ConfigName
}
