package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowEnv struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowEnv) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleWorkflowEnv) Lint(w *workflow.Workflow, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if w.Env == nil || len(w.Env) == 0 {
		return
	}

	if r.Value == "uppercase-underscores" {
		reName := regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

		for envName := range w.Env {
			m := reName.MatchString(envName)
			if !m {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' env '%s' must be alphanumeric uppercase and underscore only", w.DisplayName, envName))
			}
		}
	}

	return
}

func (r RuleWorkflowEnv) GetConfigName() string {
	return r.ConfigName
}
