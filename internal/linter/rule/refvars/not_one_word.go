package refvars

import (
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// NotOneWord checks for variable references that are single-word or single-level, e.g. `${{ something }}` instead of
// `${{ inputs.something }}`.
// Only the values `true` and `false` are permitted in this form; all other variables are considered invalid.
type NotOneWord struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r NotOneWord) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "referenced_variables_in_workflows__not_one_word"
	case rule.DotGithubFileTypeAction:
		return "referenced_variables_in_actions__not_one_word"
	default:
		return "referenced_variables_in_*__not_one_word"
	}
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r NotOneWord) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r NotOneWord) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errValueNotBool
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r NotOneWord) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	confValue, confIsBool := conf.(bool)
	if !confIsBool {
		return false, errValueNotBool
	}

	if file.GetType() != rule.DotGithubFileTypeAction &&
		file.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	if !confValue {
		return true, nil
	}

	refVarRegexp := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\-_]+)[ ]*}}`)

	compliant := true

	if file.GetType() == rule.DotGithubFileTypeAction {
		actionInstance, ok := file.(*action.Action)
		if !ok {
			return false, errFileInvalidType
		}

		found := refVarRegexp.FindAllSubmatch(actionInstance.Raw, -1)
		for _, variableReference := range found {
			if string(variableReference[1]) != "false" && string(variableReference[1]) != "true" {
				chErrors <- glitch.Glitch{
					Path:     actionInstance.Path,
					Name:     actionInstance.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  fmt.Sprintf("calls a variable '%s' that is invalid", string(variableReference[1])),
					RuleName: r.ConfigName(rule.DotGithubFileTypeAction),
				}

				compliant = false
			}
		}
	}

	if file.GetType() == rule.DotGithubFileTypeWorkflow {
		workflowInstance, ok := file.(*workflow.Workflow)
		if !ok {
			return false, errFileInvalidType
		}

		found := refVarRegexp.FindAllSubmatch(workflowInstance.Raw, -1)
		for _, variableReference := range found {
			if string(variableReference[1]) != "false" && string(variableReference[1]) != "true" {
				chErrors <- glitch.Glitch{
					Path:     workflowInstance.Path,
					Name:     workflowInstance.DisplayName,
					Type:     rule.DotGithubFileTypeWorkflow,
					ErrText:  fmt.Sprintf("calls a variable '%s' that is invalid", string(variableReference[1])),
					RuleName: r.ConfigName(rule.DotGithubFileTypeWorkflow),
				}

				compliant = false
			}
		}
	}

	return compliant, nil
}
