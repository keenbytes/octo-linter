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

// NotOneWord checks for variable references that are single-word or single-level, e.g. `${{ something }}` instead of `${{ inputs.something }}`.
// Only the values `true` and `false` are permitted in this form; all other variables are considered invalid.
type NotOneWord struct {
}

func (r NotOneWord) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "referenced_variables_in_workflows__not_one_word"
	case rule.DotGithubFileTypeAction:
		return "referenced_variables_in_actions__not_one_word"
	default:
		return "referenced_variables_in_*__not_one_word"
	}
}

func (r NotOneWord) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r NotOneWord) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r NotOneWord) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	if !conf.(bool) {
		return
	}

	re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\-_]+)[ ]*}}`)

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)

		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				chErrors <- fmt.Sprintf("action '%s' calls a variable '%s' that is invalid", a.DirName, string(f[1]))
				compliant = false
			}
		}
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)

		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that is invalid", w.FileName, string(f[1]))
				compliant = false
			}
		}
	}

	return
}
