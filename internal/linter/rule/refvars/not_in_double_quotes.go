package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// NotInDoubleQuotes scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.
type NotInDoubleQuotes struct {
}

func (r NotInDoubleQuotes) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "referenced_variables_in_workflows__not_in_double_quotes"
	case rule.DotGithubFileTypeAction:
		return "referenced_variables_in_actions__not_in_double_quotes"
	default:
		return "referenced_variables_in_*__not_in_double_quotes"
	}
}

func (r NotInDoubleQuotes) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r NotInDoubleQuotes) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r NotInDoubleQuotes) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	if !conf.(bool) {
		return
	}

	re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\-_.]+)[ ]*}}\"`)

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)

		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			chErrors <- fmt.Sprintf("action '%s' calls a variable '%s' that is in double quotes", a.DirName, string(f[1]))
			compliant = false
		}
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)

		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that is in double quotes", w.FileName, string(f[1]))
			compliant = false
		}
	}

	return
}
