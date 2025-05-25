package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

// RuleWorkflowCalledVariableNotInDoubleQuote scans for all variable references enclosed in double quotes.
// It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.
type RuleWorkflowCalledVariableNotInDoubleQuote struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) Validate() error {
	return nil
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if r.Value {
		re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls a variable '%s' that is in double quotes", w.FileName, string(f[1])), chWarnings, chErrors)
			compliant = false
		}
	}

	return
}

func (r RuleWorkflowCalledVariableNotInDoubleQuote) GetConfigName() string {
	return r.ConfigName
}
