package workflowrunners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// NotLatest checks whether 'runs-on' does not contain the 'latest' string. In some case, runner version (image) should be frozen, instead of using the latest.
type NotLatest struct {
}

func (r NotLatest) ConfigName(int) string {
	return "workflow_runners__not_latest"
}

func (r NotLatest) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r NotLatest) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r NotLatest) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !conf.(bool) || w.Jobs == nil || len(w.Jobs) == 0 {
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
				chErrors <- fmt.Sprintf("workflow '%s' job '%s' should not use 'latest' in 'runs-on' field", w.FileName, jobName)
			}
		}

		runsOnList, ok := job.RunsOn.([]interface{})
		if ok {
			for _, runsOn := range runsOnList {
				runsOnStr, ok2 := runsOn.(string)
				if ok2 && strings.Contains(runsOnStr, "latest") {
					compliant = false
					chErrors <- fmt.Sprintf("workflow '%s' job '%s' should not use 'latest' in 'runs-on' field", w.FileName, jobName)
				}
			}
		}
	}

	return
}
