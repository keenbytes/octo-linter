package dependencies

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// WorkflowNeedsWithExistingJobs checks if `needs` field references existing jobs.
type WorkflowNeedsWithExistingJobs struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r WorkflowNeedsWithExistingJobs) ConfigName(int) string {
	return "dependencies__workflow_needs_field_must_contain_already_existing_jobs"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r WorkflowNeedsWithExistingJobs) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r WorkflowNeedsWithExistingJobs) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r WorkflowNeedsWithExistingJobs) Lint(
	conf interface{},
	f dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeWorkflow || !conf.(bool) {
		return true, nil
	}

	w := f.(*workflow.Workflow)

	if len(w.Jobs) == 0 {
		return true, nil
	}

	compliant := true

	for jobName, job := range w.Jobs {
		if job.Needs != nil {
			needsStr, ok := job.Needs.(string)
			if ok {
				if w.Jobs[needsStr] == nil {
					compliant = false

					chErrors <- glitch.Glitch{
						Path:     w.Path,
						Name:     w.DisplayName,
						Type:     rule.DotGithubFileTypeWorkflow,
						ErrText:  fmt.Sprintf("job '%s' has non-existing job '%s' in 'needs' field", jobName, needsStr),
						RuleName: r.ConfigName(0),
					}
				}
			}

			needsList, ok := job.Needs.([]interface{})
			if ok {
				for _, neededJob := range needsList {
					if w.Jobs[neededJob.(string)] == nil {
						compliant = false

						chErrors <- glitch.Glitch{
							Path:     w.Path,
							Name:     w.DisplayName,
							Type:     rule.DotGithubFileTypeWorkflow,
							ErrText:  fmt.Sprintf("job '%s' has non-existing job '%s' in 'needs' field", jobName, neededJob.(string)),
							RuleName: r.ConfigName(0),
						}
					}
				}
			}
		}
	}

	return compliant, nil
}
