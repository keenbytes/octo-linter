package usedactions

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/step"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// Source checks if referenced action (in `uses`) in steps has valid path.
// This rule can be configured to allow local actions, external actions, or both.
type Source struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r Source) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "used_actions_in_workflow_job_steps__source"
	case rule.DotGithubFileTypeAction:
		return "used_actions_in_action_steps__source"
	default:
		return "used_actions_in_*_steps__source"
	}
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r Source) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r Source) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != ValueLocalOnly && val != ValueLocalOrExternal && val != ValueExternalOnly && val != "" {
		return fmt.Errorf(
			"%s supports '%s', '%s', '%s' or empty value only",
			r.ConfigName(0),
			ValueLocalOnly,
			ValueLocalOrExternal,
			ValueExternalOnly,
		)
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r Source) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if file.GetType() != rule.DotGithubFileTypeAction &&
		file.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	confVal := conf.(string)
	if confVal == "" {
		return true, nil
	}

	reLocal := regexp.MustCompile(
		`^\.\/\.github\/actions\/([a-zA-Z0-9\-_]+|[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-_]+)$`,
	)
	reExternal := regexp.MustCompile(
		`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`,
	)

	steps := []*step.Step{}
	msgPrefix := map[int]string{}

	var (
		fileType int
		filePath string
		fileName string
	)

	if file.GetType() == rule.DotGithubFileTypeAction {
		actionInstance := file.(*action.Action)
		if len(actionInstance.Runs.Steps) == 0 {
			return true, nil
		}

		steps = actionInstance.Runs.Steps
		msgPrefix[0] = ""

		fileType = rule.DotGithubFileTypeAction
		filePath = actionInstance.Path
		fileName = actionInstance.DirName
	}

	if file.GetType() == rule.DotGithubFileTypeWorkflow {
		workflowInstance := file.(*workflow.Workflow)
		if len(workflowInstance.Jobs) == 0 {
			return true, nil
		}

		for jobName, job := range workflowInstance.Jobs {
			if len(job.Steps) == 0 {
				continue
			}

			msgPrefix[len(steps)] = fmt.Sprintf("job '%s' ", jobName)

			steps = append(steps, job.Steps...)
		}

		fileType = rule.DotGithubFileTypeWorkflow
		filePath = workflowInstance.Path
		fileName = workflowInstance.DisplayName
	}

	var errPrefix string
	if file.GetType() == rule.DotGithubFileTypeAction {
		errPrefix = msgPrefix[0]
	}

	compliant := true

	for stepIdx, step := range steps {
		newErrPrefix, ok := msgPrefix[stepIdx]
		if ok {
			errPrefix = newErrPrefix
		}

		if step.Uses == "" {
			continue
		}

		isLocal := reLocal.MatchString(step.Uses)
		isExternal := reExternal.MatchString(step.Uses)

		if confVal == ValueLocalOnly && !isLocal {
			chErrors <- glitch.Glitch{
				Path: filePath,
				Name: fileName,
				Type: fileType,
				ErrText: fmt.Sprintf(
					"%sstep %d calls action '%s' that is not a valid local path",
					errPrefix,
					stepIdx+1,
					step.Uses,
				),
				RuleName: r.ConfigName(fileType),
			}

			compliant = false
		}

		if confVal == ValueExternalOnly && !isExternal {
			chErrors <- glitch.Glitch{
				Path: filePath,
				Name: fileName,
				Type: fileType,
				ErrText: fmt.Sprintf(
					"%sstep %d calls action '%s' that is not external",
					errPrefix,
					stepIdx+1,
					step.Uses,
				),
				RuleName: r.ConfigName(fileType),
			}

			compliant = false
		}

		if confVal == ValueLocalOrExternal && !isLocal && !isExternal {
			chErrors <- glitch.Glitch{
				Path: filePath,
				Name: fileName,
				Type: fileType,
				ErrText: fmt.Sprintf(
					"%sstep %d calls action '%s' that is neither external nor local",
					errPrefix,
					stepIdx+1,
					step.Uses,
				),
				RuleName: r.ConfigName(fileType),
			}

			compliant = false
		}
	}

	return compliant, nil
}
