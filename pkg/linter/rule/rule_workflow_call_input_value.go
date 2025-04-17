package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowCallInputValue struct {
	Value      map[string]string
	ConfigName string
	LogLevel   int
	IsError    map[string]bool
}

func (r RuleWorkflowCallInputValue) Validate() error {
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

func (r RuleWorkflowCallInputValue) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
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

	re := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
	for inputName := range w.On.WorkflowCall.Inputs {
		for k, v := range r.Value {
			if k == "name" && v != "" {
				m := re.MatchString(inputName)
				if !m {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError[k], r.LogLevel, fmt.Sprintf("workflow '%s' call input '%s' %s must be lower-case and hyphens only", w.FileName, inputName, v), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleWorkflowCallInputValue) GetConfigName() string {
	return r.ConfigName
}
