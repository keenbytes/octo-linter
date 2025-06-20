package usedactions

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/step"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// Source checks if referenced action (in `uses`) in steps has valid path.
// This rule can be configured to allow local actions, external actions, or both.
type Source struct {
}

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

func (r Source) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r Source) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != "local-only" && val != "local-or-external" && val != "external-only" && val != "" {
		return fmt.Errorf("%s supports 'local-only', 'external-only', 'local-or-external' or empty value only", r.ConfigName(0))
	}

	return nil
}

func (r Source) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	confVal := conf.(string)
	if confVal == "" {
		return
	}
	
	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

	steps := []*step.Step{}
	msgPrefix := map[int]string{}

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)
		if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
			return
		}
		steps = a.Runs.Steps
		msgPrefix[0] = fmt.Sprintf("action '%s'", a.DirName)
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
			msgPrefix[len(steps)] = fmt.Sprintf("workflow '%s' job '%s'", w.FileName, jobName)
			steps = append(steps, job.Steps...)
		}
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

		if confVal == "local-only" && !isLocal {
			chErrors <- fmt.Sprintf("%s step %d calls action '%s' that is not a valid local path", errPrefix, i+1, st.Uses)
			compliant = false
		}
		if confVal == "external-only" && !isExternal {
			chErrors <- fmt.Sprintf("%s step %d calls action '%s' that is not external", errPrefix, i+1, st.Uses)
			compliant = false
		}
		if confVal == "local-or-external" && !isLocal && !isExternal {
			chErrors <- fmt.Sprintf("%s step %d calls action '%s' that is neither external nor local", errPrefix, i+1, st.Uses)
			compliant = false
		}
	}

	return
}
