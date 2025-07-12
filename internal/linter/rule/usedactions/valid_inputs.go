package usedactions

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/step"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// ValidInputs verifies that all required inputs are provided when referencing an action in a step, and that no undefined inputs are used.
type ValidInputs struct {
}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r ValidInputs) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "used_actions_in_workflow_job_steps__must_have_valid_inputs"
	case rule.DotGithubFileTypeAction:
		return "used_actions_in_action_steps__must_have_valid_inputs"
	default:
		return "used_actions_in_*_steps__must_have_valid_inputs"
	}
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r ValidInputs) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r ValidInputs) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r ValidInputs) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	if !conf.(bool) {
		return true, nil
	}

	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

	steps := []*step.Step{}
	msgPrefix := map[int]string{}

	var (
		fileType int
		filePath string
		fileName string
	)

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)
		if len(a.Runs.Steps) == 0 {
			return true, nil
		}

		steps = a.Runs.Steps
		msgPrefix[0] = ""

		fileType = rule.DotGithubFileTypeAction
		filePath = a.Path
		fileName = a.DirName
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)
		if len(w.Jobs) == 0 {
			return true, nil
		}

		for jobName, job := range w.Jobs {
			if len(job.Steps) == 0 {
				continue
			}

			msgPrefix[len(steps)] = fmt.Sprintf("job '%s'", jobName)

			steps = append(steps, job.Steps...)
		}

		fileType = rule.DotGithubFileTypeWorkflow
		filePath = w.Path
		fileName = w.DisplayName
	}

	var errPrefix string
	if f.GetType() == rule.DotGithubFileTypeAction {
		errPrefix = msgPrefix[0]
	}

	compliant := true

	for i, st := range steps {
		newErrPrefix, ok := msgPrefix[i]
		if ok {
			errPrefix = newErrPrefix
		}

		if st.Uses == "" {
			continue
		}

		isLocal := reLocal.MatchString(st.Uses)
		isExternal := reExternal.MatchString(st.Uses)

		var action *action.Action

		if isLocal {
			actionName := strings.ReplaceAll(st.Uses, "./.github/actions/", "")
			action = d.GetAction(actionName)
		}

		if isExternal {
			action = d.GetExternalAction(st.Uses)
		}

		if action == nil {
			continue
		}

		if action.Inputs != nil {
			for daInputName, daInput := range action.Inputs {
				if daInput.Required {
					if st.With == nil || st.With[daInputName] == "" {
						chErrors <- glitch.Glitch{
							Path:     filePath,
							Name:     fileName,
							Type:     fileType,
							ErrText:  fmt.Sprintf("%sstep %d called action requires input '%s'", errPrefix, i+1, daInputName),
							RuleName: r.ConfigName(fileType),
						}

						compliant = false
					}
				}
			}
		}

		if st.With != nil {
			for usedInput := range st.With {
				if action.Inputs == nil || action.Inputs[usedInput] == nil {
					chErrors <- glitch.Glitch{
						Path:     filePath,
						Name:     fileName,
						Type:     fileType,
						ErrText:  fmt.Sprintf("%sstep %d called action non-existing input '%s'", errPrefix, i+1, usedInput),
						RuleName: r.ConfigName(fileType),
					}

					compliant = false
				}
			}
		}
	}

	return compliant, nil
}
