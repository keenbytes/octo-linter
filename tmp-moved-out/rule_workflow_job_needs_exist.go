package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowJobNeedsExist checks if 'needs' references existing jobs.
type RuleWorkflowJobNeedsExist struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowJobNeedsExist) Validate() error {
	return nil
}

func (r RuleWorkflowJobNeedsExist) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.Needs != nil {
			needsStr, ok := job.Needs.(string)
			if ok {
				if w.Jobs[needsStr] == nil {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' has non-existing job '%s' in 'needs' field", w.FileName, jobName, needsStr), chWarnings, chErrors)
				}
			}

			needsList, ok := job.Needs.([]interface{})
			if ok {
				for _, neededJob := range needsList {
					if w.Jobs[neededJob.(string)] == nil {
						compliant = false
						printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' has non-existing job '%s' in 'needs' field", w.FileName, jobName, neededJob.(string)), chWarnings, chErrors)
					}
				}
			}
		}
	}

	return
}

func (r RuleWorkflowJobNeedsExist) GetConfigName() string {
	return r.ConfigName
}
