package rule

import (
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowJobEnv checks whether workflow job environment variable names follow the specified naming convention.
// Currently, only 'uppercase-underscores' is supported, meaning variable names may contain uppercase letters, numbers, and underscores only.
type RuleWorkflowJobEnv struct {
	Value      string
	ConfigName string
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
					printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' job '%s' env '%s' must be alphanumeric uppercase and underscore only", w.DisplayName, jobName, envName), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleWorkflowJobEnv) GetConfigName() string {
	return r.ConfigName
}
