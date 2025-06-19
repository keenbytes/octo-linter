package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// NotOneWord_InAction checks for variable references that are single-word or single-level, e.g. `${{ something }}` instead of `${{ inputs.something }}`.
// Only the values `true` and `false` are permitted in this form; all other variables are considered invalid.
type NotOneWord_InAction struct {
}

func (r NotOneWord_InAction) ConfigName(int) string {
	return "referenced_variables_in_actions__not_one_word"
}

func (r NotOneWord_InAction) FileType() int {
	return rule.DotGithubFileTypeAction
}

func (r NotOneWord_InAction) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r NotOneWord_InAction) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return
	}
	a := f.(*action.Action)

	re := regexp.MustCompile(`\${{[ ]*([a-zA-Z0-9\\-_]+)[ ]*}}`)
	found := re.FindAllSubmatch(a.Raw, -1)
	for _, f := range found {
		if string(f[1]) != "false" && string(f[1]) != "true" {
			chErrors <- fmt.Sprintf("action '%s' calls a variable '%s' that is invalid", a.DirName, string(f[1]))
			compliant = false
		}
	}

	return
}
