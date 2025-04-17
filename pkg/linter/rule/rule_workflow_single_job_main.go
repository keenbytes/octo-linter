package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowSingleJobMain struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowSingleJobMain) Validate() error {
	return nil
}

func (r RuleWorkflowSingleJobMain) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil {
		return
	}

	if len(w.Jobs) == 1 {
		// there's only one
		for jobName := range w.Jobs {
			if jobName != "main" {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' has only one job and it should be called 'main'", w.FileName), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleWorkflowSingleJobMain) GetConfigName() string {
	return r.ConfigName
}
