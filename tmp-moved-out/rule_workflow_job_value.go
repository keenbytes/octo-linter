package rule

import (
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowJobValue checks if workflow job fields follow specified naming convention, for example
// if 'name' is 'lowercase-hyphens'.
type RuleWorkflowJobValue struct {
	Value      map[string]string
	ConfigName string
	IsError    map[string]bool
}

func (r RuleWorkflowJobValue) Validate() error {
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

func (r RuleWorkflowJobValue) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	reName := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)

	for jobName := range w.Jobs {
		for k, v := range r.Value {
			if k == "name" && v != "" {
				m := reName.MatchString(jobName)
				if !m {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError[k], fmt.Sprintf("workflow '%s' job '%s' name must be lower-case and hyphens only", w.FileName, jobName), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleWorkflowJobValue) GetConfigName() string {
	return r.ConfigName
}
