package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// NotLatest checks whether 'runs-on' does not contain the 'latest' string. In some case, runner version (image)
// should be frozen, instead of using the latest.
type NotLatest struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r NotLatest) ConfigName(int) string {
	return "workflow_runners__not_latest"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r NotLatest) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r NotLatest) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r NotLatest) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if file.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	workflowInstance := file.(*workflow.Workflow)

	if !conf.(bool) || len(workflowInstance.Jobs) == 0 {
		return true, nil
	}

	compliant := true

	for jobName, job := range workflowInstance.Jobs {
		if job.RunsOn == nil {
			continue
		}

		runsOnStr, runsOnIsString := job.RunsOn.(string)
		if runsOnIsString {
			if strings.Contains(runsOnStr, "latest") {
				compliant = false

				chErrors <- glitch.Glitch{
					Path:     workflowInstance.Path,
					Name:     workflowInstance.DisplayName,
					Type:     rule.DotGithubFileTypeWorkflow,
					ErrText:  fmt.Sprintf("job '%s' should not use 'latest' in 'runs-on' field", jobName),
					RuleName: r.ConfigName(0),
				}
			}
		}

		runsOnList, runsOnIsList := job.RunsOn.([]interface{})
		if runsOnIsList {
			for _, runsOn := range runsOnList {
				runsOnStr, ok2 := runsOn.(string)
				if ok2 && strings.Contains(runsOnStr, "latest") {
					compliant = false

					chErrors <- glitch.Glitch{
						Path:     workflowInstance.Path,
						Name:     workflowInstance.DisplayName,
						Type:     rule.DotGithubFileTypeWorkflow,
						ErrText:  fmt.Sprintf("job '%s' should not use 'latest' in 'runs-on' field", jobName),
						RuleName: r.ConfigName(0),
					}
				}
			}
		}
	}

	return compliant, nil
}
