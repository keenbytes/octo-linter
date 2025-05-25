package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/step"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

// RuleStepAction checks whether the referenced actions have valid paths.
// This rule can be configured to allow local actions, external actions, or both.
type RuleStepAction struct {
	Value      string
	ConfigName string
	IsError    bool
}

func (r RuleStepAction) Validate() error {
	if r.Value != "" {
		if r.Value != "local-only" && r.Value != "local-or-external" && r.Value != "external-only" {
			return fmt.Errorf("%s supports 'local-only', 'external-only', 'local-or-external' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleStepAction) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction && f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	if r.Value == "" {
		return
	}

	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

	steps := []*step.Step{}
	msgPrefix := map[int]string{}

	if f.GetType() == DotGithubFileTypeAction {
		a := f.(*action.Action)
		if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
			return
		}
		steps = a.Runs.Steps
		msgPrefix[0] = fmt.Sprintf("action '%s'", a.DirName)
	}
	if f.GetType() == DotGithubFileTypeWorkflow {
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
	if f.GetType() == DotGithubFileTypeAction {
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

		if r.Value == "local-only" && !isLocal {
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("%s step %d calls action '%s' that is not a valid local path", errPrefix, i+1, st.Uses), chWarnings, chErrors)
			compliant = false
		}
		if r.Value == "external-only" && !isExternal {
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("%s step %d calls action '%s' that is not external", errPrefix, i+1, st.Uses), chWarnings, chErrors)
			compliant = false
		}
		if r.Value == "local-or-external" && !isLocal && !isExternal {
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("%s step %d calls action '%s' that is neither external nor local", errPrefix, i+1, st.Uses), chWarnings, chErrors)
			compliant = false
		}
	}

	return
}

func (r RuleStepAction) GetConfigName() string {
	return r.ConfigName
}
