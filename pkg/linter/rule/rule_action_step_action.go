package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionStepAction struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionStepAction) Validate() error {
	if r.Value != "" {
		if r.Value != "local-only" && r.Value != "local-or-external" && r.Value != "external-only" {
			return fmt.Errorf("%s supports 'local-only', 'external-only', 'local-or-external' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleActionStepAction) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if r.Value == "" || a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
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
		if r.Value == "local-only" && !isLocal {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d calls action '%s' that is not a valid local path", a.DirName, i+1, step.Uses), chWarnings, chErrors)
			compliant = false
		}
		if r.Value == "external-only" && !isExternal {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d calls action '%s' that is not external", a.DirName, i+1, step.Uses), chWarnings, chErrors)
			compliant = false
		}
		if r.Value == "local-or-external" && !isLocal && !isExternal {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' step %d calls action '%s' that is neither external nor local", a.DirName, i+1, step.Uses), chWarnings, chErrors)
			compliant = false
		}
	}

	return
}

func (r RuleActionStepAction) GetConfigName() string {
	return r.ConfigName
}
