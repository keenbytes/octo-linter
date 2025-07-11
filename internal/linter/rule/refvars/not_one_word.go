package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
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

func (r NotOneWord) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	if !conf.(bool) {
		return true, nil
	}

	re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\-_]+)[ ]*}}`)

	compliant := true

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)

		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				chErrors <- glitch.Glitch{
					Path:     a.Path,
					Name:     a.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  fmt.Sprintf("calls a variable '%s' that is invalid", string(f[1])),
					RuleName: r.ConfigName(rule.DotGithubFileTypeAction),
				}
				compliant = false
			}
		}
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)

		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			if string(f[1]) != "false" && string(f[1]) != "true" {
				chErrors <- glitch.Glitch{
					Path:     w.Path,
					Name:     w.DisplayName,
					Type:     rule.DotGithubFileTypeWorkflow,
					ErrText:  fmt.Sprintf("calls a variable '%s' that is invalid", string(f[1])),
					RuleName: r.ConfigName(rule.DotGithubFileTypeWorkflow),
				}
				compliant = false
			}
		}
	}

	return compliant, nil
}
