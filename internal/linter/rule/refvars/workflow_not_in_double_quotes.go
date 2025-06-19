package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowNotInDoubleQuotes scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.
type WorkflowNotInDoubleQuotes struct {
}

func (r WorkflowNotInDoubleQuotes) ConfigName() string {
	return "referenced_variables_in_workflows__not_in_double_quotes"
}

func (r WorkflowNotInDoubleQuotes) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowNotInDoubleQuotes) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowNotInDoubleQuotes) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return
	}
	w := f.(*workflow.Workflow)

	re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
	found := re.FindAllSubmatch(w.Raw, -1)
	for _, f := range found {
		chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that is in double quotes", w.FileName, string(f[1]))
		compliant = false
	}

	return
}
