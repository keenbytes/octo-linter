package naming

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
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

func (r WorkflowSingleJobOnlyName) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
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
				chErrors <- fmt.Sprintf("workflow '%s' has only one job and it should be called '%s'", w.FileName, conf.(string))
				compliant = false
			}
		}
	}

	return
}
