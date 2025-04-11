package rule

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionStepActionExists struct {
	Value      []string
	ConfigName string
	LogLevel   int
	IsError    []bool
}

func (r RuleActionStepActionExists) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "local" && v != "external" {
				return fmt.Errorf("%s can only contain 'local' and/or 'external'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionStepActionExists) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if len(r.Value) == 0 || a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
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
		if isLocal {
			actionName := strings.Replace(step.Uses, "./.github/actions/", "", -1)
			action := d.GetAction(actionName)
			if action == nil {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("action '%s' step %d calls non-existing local action '%s'", a.DirName, i+1, actionName), chWarnings, chErrors)
			}
		}
		if isExternal {
			action := d.GetExternalAction(step.Uses)
			if action == nil {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError[i], r.LogLevel, fmt.Sprintf("action '%s' step %d calls non-existing external action '%s'", a.DirName, i+1, step.Uses), chWarnings, chErrors)
			}
		}
	}

	return
}

func (r RuleActionStepActionExists) GetConfigName() string {
	return r.ConfigName
}
