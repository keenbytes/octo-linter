package naming

import (
	"fmt"

	"octo-linter/internal/dotgithub"
	"octo-linter/internal/linter/glitch"
	"octo-linter/internal/linter/rule"
	"octo-linter/internal/workflow"
)

// WorkflowSingleJobOnlyName checks if a workflow has only one job, this should be its name.
type WorkflowSingleJobOnlyName struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r WorkflowSingleJobOnlyName) ConfigName(int) string {
	return "filenames__workflow_filename_base_format"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r WorkflowSingleJobOnlyName) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r WorkflowSingleJobOnlyName) Validate(conf interface{}) error {
	_, ok := conf.(string)
	if !ok {
		return errValueNotString
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r WorkflowSingleJobOnlyName) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	confValue, confIsString := conf.(string)
	if !confIsString {
		return false, errValueNotString
	}

	if file.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	workflowInstance, ok := file.(*workflow.Workflow)
	if !ok {
		return false, errFileInvalidType
	}

	if confValue == "" || workflowInstance.Jobs == nil {
		return true, nil
	}

	if len(workflowInstance.Jobs) != 1 {
		return true, nil
	}

	for jobName := range workflowInstance.Jobs {
		if jobName != confValue {
			chErrors <- glitch.Glitch{
				Path:     workflowInstance.Path,
				Name:     workflowInstance.DisplayName,
				Type:     rule.DotGithubFileTypeWorkflow,
				ErrText:  fmt.Sprintf("has only one job and it should be called '%s'", confValue),
				RuleName: r.ConfigName(0),
			}

			return false, nil
		}
	}

	return true, nil
}
