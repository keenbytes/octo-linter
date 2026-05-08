// Package refvars contains rules checking variables referenced in action or workflow steps, eg. ${{ var }}.
package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"octo-linter/internal/action"
	"octo-linter/internal/linter/glitch"
	"octo-linter/internal/linter/rule"
	"octo-linter/internal/workflow"
)

var (
	errFileInvalidType = errors.New("file is of invalid type")
	errValueNotBool    = errors.New("value should be bool")
)

var (
	regexpReferenceInDoubleQuote = regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\-_.]+)[ ]*}}\"`)
	regexpReference              = regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\-_]+)[ ]*}}`)
)

func processActionForRegexp(
	ruleConfigName string,
	actionInstance *action.Action,
	regexpToMatch *regexp.Regexp,
	chErrors chan<- glitch.Glitch,
	errorText string,
) bool {
	foundNotCompliant := false

	found := regexpToMatch.FindAllSubmatch(actionInstance.Raw, -1)
	for _, ref := range found {
		chErrors <- glitch.Glitch{
			Path:     actionInstance.Path,
			Name:     actionInstance.DirName,
			Type:     rule.DotGithubFileTypeAction,
			ErrText:  fmt.Sprintf("calls a variable '%s' that is %s", string(ref[1]), errorText),
			RuleName: ruleConfigName,
		}

		foundNotCompliant = true
	}

	return foundNotCompliant
}

func processWorkflowForRegexp(
	ruleConfigName string,
	workflowInstance *workflow.Workflow,
	regexpToMatch *regexp.Regexp,
	chErrors chan<- glitch.Glitch,
	errorText string,
) bool {
	foundNotCompliant := false

	found := regexpToMatch.FindAllSubmatch(workflowInstance.Raw, -1)
	for _, ref := range found {
		chErrors <- glitch.Glitch{
			Path:     workflowInstance.Path,
			Name:     workflowInstance.DisplayName,
			Type:     rule.DotGithubFileTypeWorkflow,
			ErrText:  fmt.Sprintf("calls a variable '%s' that is %s", string(ref[1]), errorText),
			RuleName: ruleConfigName,
		}

		foundNotCompliant = true
	}

	return foundNotCompliant
}
