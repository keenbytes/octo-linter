package rule

import (
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleActionCalledVariable verifies that referenced variables such as 'env', 'var', and 'secret' follow the defined casing rule.
// Currently, only 'uppercase-underscores' is supported, meaning variables must be fully uppercase and may include underscores.
type RuleWorkflowCalledVariable struct {
	Value      string
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowCalledVariable) Validate() error {
	if r.Value != "" {
		if r.Value != "uppercase-underscores" {
			return fmt.Errorf("%s supports 'uppercase-underscores' or empty value only", r.ConfigName)
		}
	}
	return nil
}

func (r RuleWorkflowCalledVariable) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if r.Value == "uppercase-underscores" {
		reName := regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

		varTypes := []string{"env", "vars", "secrets"}
		for _, v := range varTypes {
			re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
			found := re.FindAllSubmatch(w.Raw, -1)
			for _, f := range found {
				m := reName.MatchString(string(f[1]))
				if !m {
					printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls a variable '%s' that must be alphanumeric uppercase and underscore only", w.FileName, string(f[1])), chWarnings, chErrors)
					compliant = false
				}
			}
		}
	}

	return
}

func (r RuleWorkflowCalledVariable) GetConfigName() string {
	return r.ConfigName
}
