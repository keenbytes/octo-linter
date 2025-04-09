package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowCalledVariableNotInDoubleQuote struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) Validate() error {
	return nil
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) Lint(w *workflow.Workflow, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if r.Value {
		re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' calls a variable '%s' that is in double quotes", w.FileName, string(f[1])))
			compliant = false
		}
	}

	return
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) GetConfigName() string {
	return r.ConfigName
}
