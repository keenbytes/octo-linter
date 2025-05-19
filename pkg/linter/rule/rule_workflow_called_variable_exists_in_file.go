package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowCalledVariableExistsInFile struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowCalledVariableExistsInFile) Validate() error {
	return nil
}

func (r RuleWorkflowCalledVariableExistsInFile) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if !r.Value {
		return
	}

	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	varTypes := []string{"vars", "secrets"}
	for _, v := range varTypes {
		re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			if v == "vars" && len(d.Vars) > 0 && !d.IsVarExist(string(f[1])) {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls a variable '%s' that does not exist in the vars file", w.FileName, string(f[1])), chWarnings, chErrors)
				compliant = false
			}

			if v == "secrets" && len(d.Secrets) > 0 && !d.IsSecretExist(string(f[1])) {
				printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls a secret '%s' that does not exist in the secrets file", w.FileName, string(f[1])), chWarnings, chErrors)
				compliant = false
			}
		}
	}

	return
}

func (r RuleWorkflowCalledVariableExistsInFile) GetConfigName() string {
	return r.ConfigName
}
