package required

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// Workflow checks if workflow has `runs-on` or `uses` field. At least of them must be defined.
type WorkflowUsesOrRunsOn struct {
	Field string
}

func (r WorkflowUsesOrRunsOn) ConfigName(int) string {
	return "required_fields__workflow_requires_uses_or_runs_on_required"
}

func (r WorkflowUsesOrRunsOn) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowUsesOrRunsOn) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowUsesOrRunsOn) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	err = r.Validate(conf)
	if err != nil {
		return
	}

	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !conf.(bool) || w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.RunsOn == nil && job.Uses == "" {
			chErrors <- glitch.Glitch{
				Path:     w.Path,
				Name:     w.DisplayName,
				Type:     rule.DotGithubFileTypeWorkflow,
				ErrText:  fmt.Sprintf("job '%s' should have either 'uses' or 'runs-on' field", jobName),
				RuleName: r.ConfigName(0),
			}
			compliant = false
		}

		runsOnStr, ok := job.RunsOn.(string)
		if ok {
			if job.Uses == "" && runsOnStr == "" {
				chErrors <- glitch.Glitch{
					Path:     w.Path,
					Name:     w.DisplayName,
					Type:     rule.DotGithubFileTypeWorkflow,
					ErrText:  fmt.Sprintf("job '%s' should have either 'uses' or 'runs-on' field", jobName),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	}

	return
}
