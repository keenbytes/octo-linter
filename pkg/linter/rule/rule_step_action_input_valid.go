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

type RuleStepActionInputValid struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleStepActionInputValid) Validate() error {
	return nil
}

func (r RuleStepActionInputValid) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction && f.GetType() != DotGithubFileTypeWorkflow {
		return
	}

	if !r.Value {
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
	
		var action *action.Action
		if isLocal {
			actionName := strings.Replace(st.Uses, "./.github/actions/", "", -1)
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
						printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("%s step %d called action requires input '%s'", errPrefix, i+1, daInputName), chWarnings, chErrors)
						compliant = false
					}
				}
			}
		}
		if st.With != nil {
			for usedInput := range st.With {
				if action.Inputs == nil || action.Inputs[usedInput] == nil {
					printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("%s step %d called action non-existing input '%s'", errPrefix, i+1, usedInput), chWarnings, chErrors)
					compliant = false
				}
			}
		}
	}

	return
}

func (r RuleStepActionInputValid) GetConfigName() string {
	return r.ConfigName
}
