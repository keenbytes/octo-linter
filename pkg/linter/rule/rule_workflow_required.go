package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowRequired struct {
	Value      []string
	ConfigName string
	LogLevel   int
	IsError    []bool
}

func (r RuleWorkflowRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "name" {
				return fmt.Errorf("%s can only contain values of 'name'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleWorkflowRequired) Lint(w *workflow.Workflow, d *dotgithub.DotGithub) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	for i, v := range r.Value {
		if v == "name" && w.Name == "" {
			printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("workflow '%s' does not have a required %s", w.DisplayName, v))
		}
	}

	return true, nil
}

func (r RuleWorkflowRequired) GetConfigName() string {
	return r.ConfigName
}
