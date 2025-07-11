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
	Field int
}

const (
	_ = iota
	WorkflowFieldWorkflow
	WorkflowFieldDispatchInput
	WorkflowFieldCallInput
)

func (r Workflow) ConfigName(int) string {
	switch r.Field {
	case WorkflowFieldWorkflow:
		return "required_fields__workflow_requires"
	case WorkflowFieldDispatchInput:
		return "required_fields__workflow_dispatch_input_requires"
	case WorkflowFieldCallInput:
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
		case WorkflowFieldWorkflow:
			if field != "name" {
				return fmt.Errorf("value can contain only 'name'")
			}
		case WorkflowFieldDispatchInput, WorkflowFieldCallInput:
			if field != "description" {
				return fmt.Errorf("value can contain only 'description'")
			}
		}
	}

	return nil
}

func (r Workflow) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	w := f.(*workflow.Workflow)

	compliant := true

	confInterfaces := conf.([]interface{})

	switch r.Field {
	case WorkflowFieldWorkflow:
		for _, field := range confInterfaces {
			if field.(string) == "name" && w.Name == "" {
				chErrors <- glitch.Glitch{
					Path:     w.Path,
					Name:     w.DisplayName,
					Type:     rule.DotGithubFileTypeWorkflow,
					ErrText:  fmt.Sprintf("does not have a required %s", field.(string)),
					RuleName: r.ConfigName(0),
				}

				compliant = false
			}
		}

	case WorkflowFieldDispatchInput:
		if w.On == nil || w.On.WorkflowDispatch == nil || len(w.On.WorkflowDispatch.Inputs) == 0 {
			return true, nil
		}

		for inputName, input := range w.On.WorkflowDispatch.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- glitch.Glitch{
						Path:     w.Path,
						Name:     w.DisplayName,
						Type:     rule.DotGithubFileTypeWorkflow,
						ErrText:  fmt.Sprintf("dispatch input '%s' does not have a required %s", inputName, field.(string)),
						RuleName: r.ConfigName(0),
					}

					compliant = false
				}
			}
		}
	case WorkflowFieldCallInput:
		if w.On == nil || w.On.WorkflowCall == nil || len(w.On.WorkflowCall.Inputs) == 0 {
			return true, nil
		}

		for inputName, input := range w.On.WorkflowCall.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- glitch.Glitch{
						Path:     w.Path,
						Name:     w.DisplayName,
						Type:     rule.DotGithubFileTypeWorkflow,
						ErrText:  fmt.Sprintf("call input '%s' does not have a required %s", inputName, field.(string)),
						RuleName: r.ConfigName(0),
					}

					compliant = false
				}
			}
		}
	}

	return compliant, nil
}
