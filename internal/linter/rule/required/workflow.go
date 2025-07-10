package required

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// Workflow checks if required fields within workflow are defined
type Workflow struct {
	Field string
}

func (r Workflow) ConfigName(int) string {
	switch r.Field {
	case "workflow":
		return "required_fields__workflow_requires"
	case "dispatch_input":
		return "required_fields__workflow_dispatch_input_requires"
	case "call_input":
		return "required_fields__workflow_call_input_requires"
	default:
		return "required_fields__workflown_*_requires"
	}
}

func (r Workflow) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r Workflow) Validate(conf interface{}) error {
	vals, ok := conf.([]interface{})
	if !ok {
		return errors.New("value should be []string")
	}

	for _, v := range vals {
		field, ok := v.(string)
		if !ok {
			return errors.New("value should be []string")
		}

		switch r.Field {
		case "workflow":
			if field != "name" {
				return fmt.Errorf("value can contain only 'name'")
			}
		case "dispatch_input", "call_input":
			if field != "description" {
				return fmt.Errorf("value can contain only 'description'")
			}
		default:
			// nothing
		}
	}

	return nil
}

func (r Workflow) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	confInterfaces := conf.([]interface{})
	switch r.Field {
	case "workflow":
		for _, field := range confInterfaces {
			if field.(string) == "name" && w.Name == "" {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText: fmt.Sprintf("does not have a required %s", field.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}

	case "dispatch_input":
		if w.On == nil || w.On.WorkflowDispatch == nil || w.On.WorkflowDispatch.Inputs == nil || len(w.On.WorkflowDispatch.Inputs) == 0 {
			return
		}

		for inputName, input := range w.On.WorkflowDispatch.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- glitch.Glitch{
						Path: w.Path,
						Name: w.DisplayName,
						Type: rule.DotGithubFileTypeWorkflow,
						ErrText: fmt.Sprintf("dispatch input '%s' does not have a required %s", inputName, field.(string)),
						RuleName: r.ConfigName(0),
					}
					compliant = false
				}
			}
		}
	case "call_input":
		if w.On == nil || w.On.WorkflowCall == nil || w.On.WorkflowCall.Inputs == nil || len(w.On.WorkflowCall.Inputs) == 0 {
			return
		}

		for inputName, input := range w.On.WorkflowCall.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- glitch.Glitch{
						Path: w.Path,
						Name: w.DisplayName,
						Type: rule.DotGithubFileTypeWorkflow,
						ErrText: fmt.Sprintf("call input '%s' does not have a required %s", inputName, field.(string)),
						RuleName: r.ConfigName(0),
					}
					compliant = false
				}
			}
		}
	default:
		// do nothing
	}

	return
}
