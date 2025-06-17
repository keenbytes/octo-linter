package rule

import (
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowRunsOnNotLatest checks whether 'runs-on' does not contain the 'latest' string.
// In some case, runner version (image) should be frozen, instead of using the latest.
type RuleWorkflowRunsOnNotLatest struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowRunsOnNotLatest) Validate() error {
	return nil
}

func (r RuleWorkflowRunsOnNotLatest) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.RunsOn == nil {
			continue
		}

		runsOnStr, ok := job.RunsOn.(string)
		if ok {
			if strings.Contains(runsOnStr, "latest") {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' should not use 'latest' in 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
			}
		}

		runsOnList, ok := job.RunsOn.([]interface{})
		if ok {
			for _, runsOn := range runsOnList {
				runsOnStr, ok2 := runsOn.(string)
				if ok2 && strings.Contains(runsOnStr, "latest") {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' should not use 'latest' in 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleWorkflowRunsOnNotLatest) GetConfigName() string {
	return r.ConfigName
}
