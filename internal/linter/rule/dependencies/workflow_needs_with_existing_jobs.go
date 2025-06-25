package dependencies

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// WorkflowNeedsWithExistingJobs checks if `needs` field references existing jobs.
type WorkflowNeedsWithExistingJobs struct {
}

func (r WorkflowNeedsWithExistingJobs) ConfigName(int) string {
	return "dependencies__workflow_needs_field_must_contain_already_existing_jobs"
}

func (r WorkflowNeedsWithExistingJobs) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowNeedsWithExistingJobs) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowNeedsWithExistingJobs) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow || !conf.(bool) {
		return
	}
	w := f.(*workflow.Workflow)

	if w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.Needs != nil {
			needsStr, ok := job.Needs.(string)
			if ok {
				if w.Jobs[needsStr] == nil {
					compliant = false
					chErrors <- fmt.Sprintf("workflow '%s' job '%s' has non-existing job '%s' in 'needs' field", w.FileName, jobName, needsStr)
				}
			}

			needsList, ok := job.Needs.([]interface{})
			if ok {
				for _, neededJob := range needsList {
					if w.Jobs[neededJob.(string)] == nil {
						compliant = false
						chErrors <- fmt.Sprintf("workflow '%s' job '%s' has non-existing job '%s' in 'needs' field", w.FileName, jobName, neededJob.(string))
					}
				}
			}
		}
	}

	return
}
