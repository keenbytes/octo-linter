package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowRequiredUsesOrRunsOn checks if workflow has 'runs-on' or 'uses' field. At least of them
// must be defined.
type RuleWorkflowRequiredUsesOrRunsOn struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowRequiredUsesOrRunsOn) Validate() error {
	return nil
}

func (r RuleWorkflowRequiredUsesOrRunsOn) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.RunsOn == nil && job.Uses == "" {
			compliant = false
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' should have either 'uses' or 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
		}

		runsOnStr, ok := job.RunsOn.(string)
		if ok {
			if job.Uses == "" && runsOnStr == "" {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' should have either 'uses' or 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
			}
		}
	}

	return
}

func (r RuleWorkflowRequiredUsesOrRunsOn) GetConfigName() string {
	return r.ConfigName
}
