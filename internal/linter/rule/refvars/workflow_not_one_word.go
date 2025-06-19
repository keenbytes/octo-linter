package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowNotOneWord checks for variable references that are single-word or single-level, e.g. `${{ something }}` instead of `${{ inputs.something }}`.
// Only the values `true` and `false` are permitted in this form; all other variables are considered invalid.
type WorkflowNotOneWord struct {
}

func (r WorkflowNotOneWord) ConfigName() string {
	return "referenced_variables_in_workflows__not_one_word"
}

func (r WorkflowNotOneWord) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowNotOneWord) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowNotOneWord) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return
	}
	w := f.(*workflow.Workflow)

	re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\\-_]+)[ ]*}}`)
	found := re.FindAllSubmatch(w.Raw, -1)
	for _, f := range found {
		if string(f[1]) != "false" && string(f[1]) != "true" {
			chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that is invalid", w.FileName, string(f[1]))
			compliant = false
		}
	}

	return
}
