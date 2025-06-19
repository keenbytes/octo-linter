package dependencies

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowReferencedVariableExistsInFile checks if called variables and secrets exist.
// This rule requires a list of variables and secrets to be checked against.
type WorkflowReferencedVariableExistsInFile struct {
}

func (r WorkflowReferencedVariableExistsInFile) ConfigName(int) string {
	return "dependencies__workflow_referenced_variable_must_exists_in_attached_file"
}

func (r WorkflowReferencedVariableExistsInFile) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowReferencedVariableExistsInFile) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowReferencedVariableExistsInFile) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow || !conf.(bool) {
		return
	}
	w := f.(*workflow.Workflow)

	varTypes := []string{"vars", "secrets"}
	for _, v := range varTypes {
		re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			if v == "vars" && len(d.Vars) > 0 && !d.IsVarExist(string(f[1])) {
				chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that does not exist in the vars file", w.FileName, string(f[1]))
				compliant = false
			}

			if v == "secrets" && len(d.Secrets) > 0 && !d.IsSecretExist(string(f[1])) {
				chErrors <- fmt.Sprintf("workflow '%s' calls a secret '%s' that does not exist in the secrets file", w.FileName, string(f[1]))
				compliant = false
			}
		}
	}

	return
}
