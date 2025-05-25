package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

// RuleStepEnv checks whether step environment variable names follow the specified naming convention.
// Currently, only 'uppercase-underscores' is supported, meaning variable names may contain uppercase letters, numbers, and underscores only.
type RuleStepEnv struct {
	Value      string
	ConfigName string
	IsError    bool
}

func (r RuleStepEnv) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleStepEnv) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction && f.GetType() != DotGithubFileTypeWorkflow {
		return
	}

	var reName *regexp.Regexp
	if r.Value == "uppercase-underscores" {
		reName = regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)
	}

	// action
	if f.GetType() == DotGithubFileTypeAction {
		a := f.(*action.Action)
		if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
			return
		}

		for i, step := range a.Runs.Steps {
			if step.Env == nil || len(step.Env) == 0 {
				continue
			}
			for envName := range step.Env {
				m := reName.MatchString(envName)
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("action '%s' step %d env '%s' must be alphanumeric uppercase and underscore only", a.DirName, i, envName), chWarnings, chErrors)
					compliant = false
				}
			}
		}
	}

	// workflow
	if f.GetType() == DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)
		for jobName, job := range w.Jobs {
			for i, step := range job.Steps {
				if step.Env == nil || len(step.Env) == 0 {
					continue
				}
				for envName := range step.Env {
					m := reName.MatchString(envName)
					if !m {
						printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' step %d env '%s' must be alphanumeric uppercase and underscore only", w.FileName, jobName, i, envName), chWarnings, chErrors)
						compliant = false
					}
				}
			}
		}
	}

	return
}

func (r RuleStepEnv) GetConfigName() string {
	return r.ConfigName
}
