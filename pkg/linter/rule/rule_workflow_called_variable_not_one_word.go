package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

// RuleWorkflowCalledVariableNotOneWord checks for variable references that are single-word or single-level,
// e.g. '${{ something }}' instead of '${{ inputs.something }}'.
// Only the values 'true' and 'false' are permitted in this form; all other variables are considered invalid.
type RuleWorkflowCalledVariableNotOneWord struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowCalledVariableNotOneWord) Validate() error {
	return nil
}

func (r RuleWorkflowCalledVariableNotOneWord) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if r.Value {
		re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls a variable '%s' that is invalid", w.FileName, string(f[1])), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleWorkflowCalledVariableNotOneWord) GetConfigName() string {
	return r.ConfigName
}
