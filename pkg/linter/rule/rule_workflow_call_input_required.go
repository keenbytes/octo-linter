package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowCallInputRequired struct {
	Value      []string
	ConfigName string
	IsError    []bool
}

func (r RuleWorkflowCallInputRequired) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "description" {
				return fmt.Errorf("%s can only contain 'description'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleWorkflowCallInputRequired) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if w.On == nil || w.On.WorkflowCall == nil || w.On.WorkflowCall.Inputs == nil || len(w.On.WorkflowCall.Inputs) == 0 {
		return
	}

	for inputName, input := range w.On.WorkflowCall.Inputs {
		for i, v := range r.Value {
			if v == "description" && input.Description == "" {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError[i], fmt.Sprintf("workflow '%s' call input '%s' does not have a required %s", w.FileName, inputName, v), chWarnings, chErrors)
			}
		}
	}

	return
}

func (r RuleWorkflowCallInputRequired) GetConfigName() string {
	return r.ConfigName
}
