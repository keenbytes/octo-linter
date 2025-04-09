package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowDispatchInputRequired struct {
	Value      []string
	ConfigName string
	LogLevel   int
	IsError    []bool
}

func (r RuleWorkflowDispatchInputRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "description" {
				return fmt.Errorf("%s can only contain 'description'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleWorkflowDispatchInputRequired) Lint(w *workflow.Workflow, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if w.On == nil || w.On.WorkflowDispatch == nil || w.On.WorkflowDispatch.Inputs == nil || len(w.On.WorkflowDispatch.Inputs) == 0 {
		return
	}

	for inputName, input := range w.On.WorkflowDispatch.Inputs {
		for i, v := range r.Value {
			if v == "description" && input.Description == "" {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("workflow '%s' dispatch input '%s' does not have a required %s", w.FileName, inputName, v))
			}
		}
	}

	return
}

func (r RuleWorkflowDispatchInputRequired) GetConfigName() string {
	return r.ConfigName
}
