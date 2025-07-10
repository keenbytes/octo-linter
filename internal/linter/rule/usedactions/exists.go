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

// Exists verifies that the action referenced in a step actually exists.
type Exists struct {
}

func (r Exists) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "used_actions_in_workflow_job_steps__must_exist"
	case rule.DotGithubFileTypeAction:
		return "used_actions_in_action_steps__must_exist"
	default:
		return "used_actions_in_*_steps__must_exist"
	}
}

func (r Exists) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r Exists) Validate(conf interface{}) error {
	vals, ok := conf.([]interface{})
	if !ok {
		return errors.New("value should be []string")
	}

	for _, v := range vals {
		source, ok := v.(string)
		if !ok {
			return errors.New("value should be []string")
		}
		if source != "local" && source != "external" {
			return fmt.Errorf("value can contain only 'local' and/or 'external'")
		}
	}

	return nil
}

func (r Exists) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	var checkLocal bool
	var checkExternal bool

	valInterfaces := conf.([]interface{})
	for _, v := range valInterfaces {
		if v == "local" {
			checkLocal = true
		}
		if v == "external" {
			checkExternal = true
		}
	}

	if !checkLocal && !checkExternal {
		return
	}

	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]+){0,1}@[a-zA-Z0-9\.\-\_]+`)

	steps := []*step.Step{}
	msgPrefix := map[int]string{}

	var fileType int
	var filePath string
	var fileName string

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)
		if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
			return
		}
		steps = a.Runs.Steps
		msgPrefix[0] = ""

		fileType = rule.DotGithubFileTypeAction
		filePath = a.Path
		fileName = a.DirName
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)
		if w.Jobs == nil || len(w.Jobs) == 0 {
			return
		}
		for jobName, job := range w.Jobs {
			if job.Steps == nil || len(job.Steps) == 0 {
				continue
			}
			msgPrefix[len(steps)] = fmt.Sprintf("job '%s' ", jobName)
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

		if checkLocal && isLocal {
			actionName := strings.Replace(st.Uses, "./.github/actions/", "", -1)
			action := d.GetAction(actionName)
			if action == nil {
				compliant = false
				chErrors <- glitch.Glitch{
					Path:     filePath,
					Name:     fileName,
					Type:     fileType,
					ErrText:  fmt.Sprintf("%sstep %d calls non-existing local action '%s'", errPrefix, i+1, actionName),
					RuleName: r.ConfigName(fileType),
				}
			}
		}
		if checkExternal && isExternal {
			action := d.GetExternalAction(st.Uses)
			if action == nil {
				compliant = false
				chErrors <- glitch.Glitch{
					Path:     filePath,
					Name:     fileName,
					Type:     fileType,
					ErrText:  fmt.Sprintf("%sstep %d calls non-existing external action '%s'", errPrefix, i+1, st.Uses),
					RuleName: r.ConfigName(fileType),
				}
			}
		}
	}

	return
}
