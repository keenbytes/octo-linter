package naming

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// WorkflowSingleJobOnlyName checks if workflow has only one job, this should be its name.
type WorkflowSingleJobOnlyName struct {
}

func (r WorkflowSingleJobOnlyName) ConfigName(int) string {
	return "filenames__workflow_filename_base_format"
}

func (r WorkflowSingleJobOnlyName) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowSingleJobOnlyName) Validate(conf interface{}) error {
	_, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	return nil
}

func (r WorkflowSingleJobOnlyName) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if conf.(string) == "" || w.Jobs == nil {
		return
	}

	if len(w.Jobs) == 1 {
		for jobName := range w.Jobs {
			if jobName != conf.(string) {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText: fmt.Sprintf("has only one job and it should be called '%s'", conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	}

	return
}
