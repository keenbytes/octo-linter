package rule

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionStepActionInputValid struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionStepActionInputValid) Validate() error {
	return nil
}

func (r RuleActionStepActionInputValid) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if !r.Value || a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
		return
	}

	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

	for i, step := range a.Runs.Steps {
		if step.Uses == "" {
			continue
		}

		isLocal := reLocal.MatchString(step.Uses)
		isExternal := reExternal.MatchString(step.Uses)
		var action *action.Action
		if isLocal {
			actionName := strings.Replace(step.Uses, "./.github/actions/", "", -1)
			action = d.GetAction(actionName)
		}
		if isExternal {
			action = d.GetExternalAction(step.Uses)
		}
		if action == nil {
			continue
		}

		if action.Inputs != nil {
			for daInputName, daInput := range action.Inputs {
				if daInput.Required {
					if step.With == nil || step.With[daInputName] == "" {
						printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d called action requires input '%s'", a.DirName, i+1, daInputName), chWarnings, chErrors)
						compliant = false
					}
				}
			}
		}
		if step.With != nil {
			for usedInput := range step.With {
				if action.Inputs == nil || action.Inputs[usedInput] == nil {
					printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d called action non-existing input '%s'", a.DirName, i+1, usedInput), chWarnings, chErrors)
					compliant = false
				}
			}
		}
	}

	return
}

func (r RuleActionStepActionInputValid) GetConfigName() string {
	return r.ConfigName
}
