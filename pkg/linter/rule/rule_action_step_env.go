package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionStepEnv struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionStepEnv) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleActionStepEnv) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
		return
	}

	if r.Value == "uppercase-underscores" {
		reName := regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

		for i, step := range a.Runs.Steps {
			if step.Env == nil || len(step.Env) == 0 {
				continue
			}
			for envName := range step.Env {
				m := reName.MatchString(envName)
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d env '%s' must be alphanumeric uppercase and underscore only", a.DirName, i, envName), chWarnings, chErrors)
					compliant = false
				}
			}
		}
	}

	return
}

func (r RuleActionStepEnv) GetConfigName() string {
	return r.ConfigName
}
