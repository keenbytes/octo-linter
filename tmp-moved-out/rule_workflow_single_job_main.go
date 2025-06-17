package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowSingleJobMain checks if workflow's only job is called 'main' - just for naming
// consistency.
type RuleWorkflowSingleJobMain struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowSingleJobMain) Validate() error {
	return nil
}

func (r RuleWorkflowSingleJobMain) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil {
		return
	}

	if len(w.Jobs) == 1 {
		// there's only one
		for jobName := range w.Jobs {
			if jobName != "main" {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' has only one job and it should be called 'main'", w.FileName), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleWorkflowSingleJobMain) GetConfigName() string {
	return r.ConfigName
}
