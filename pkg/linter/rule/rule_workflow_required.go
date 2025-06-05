package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowRequired checks whether the specified workflow fields are present, e.g. 'name'.
type RuleWorkflowRequired struct {
	Value      []string
	ConfigName string
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

func (r RuleWorkflowRequired) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	if len(r.Value) == 0 {
		return true, nil
	}

	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	for i, v := range r.Value {
		if v == "name" && w.Name == "" {
			printErrOrWarn(r.ConfigName, r.IsError[i], fmt.Sprintf("workflow '%s' does not have a required %s", w.DisplayName, v), chWarnings, chErrors)
		}
	}

	return true, nil
}

func (r RuleWorkflowRequired) GetConfigName() string {
	return r.ConfigName
}
