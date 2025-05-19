package rule

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/step"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleStepActionExists struct {
	Value      []string
	ConfigName string
	IsError    []bool
}

func (r RuleStepActionExists) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "local" && v != "external" {
				return fmt.Errorf("%s can only contain 'local' and/or 'external'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleStepActionExists) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction && f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	if len(r.Value) == 0 {
		return
	}

	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]+){0,1}@[a-zA-Z0-9\.\-\_]+`)

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

		for k, v := range r.Value {
			if v == "local" && isLocal {
				actionName := strings.Replace(st.Uses, "./.github/actions/", "", -1)
				action := d.GetAction(actionName)
				if action == nil {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError[k], fmt.Sprintf("%s step %d calls non-existing local action '%s'", errPrefix, i+1, actionName), chWarnings, chErrors)
				}
			}
			if v == "external" && isExternal {
				action := d.GetExternalAction(st.Uses)
				if action == nil {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError[k], fmt.Sprintf("%s step %d calls non-existing external action '%s'", errPrefix, i+1, st.Uses), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleStepActionExists) GetConfigName() string {
	return r.ConfigName
}
