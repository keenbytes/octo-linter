package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowJobEnv struct {
	Value      string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowJobEnv) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleWorkflowJobEnv) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	if r.Value == "uppercase-underscores" {
		reName := regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

		for jobName, job := range w.Jobs {
			if job.Env == nil || len(job.Env) == 0 {
				continue
			}
			for envName := range job.Env {
				m := reName.MatchString(envName)
			if !m {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' job '%s' env '%s' must be alphanumeric uppercase and underscore only", w.DisplayName, jobName, envName), chWarnings, chErrors)
			}
			}
		}
	}

	return
}

func (r RuleWorkflowJobEnv) GetConfigName() string {
	return r.ConfigName
}
