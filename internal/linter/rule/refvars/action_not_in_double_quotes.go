package refvars

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// ActionNotInDoubleQuotes scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.
type ActionNotInDoubleQuotes struct {
}

func (r ActionNotInDoubleQuotes) ConfigName() string {
	return "referenced_variables_in_actions__not_in_double_quotes"
}

func (r ActionNotInDoubleQuotes) FileType() int {
	return rule.DotGithubFileTypeAction
}

func (r ActionNotInDoubleQuotes) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r ActionNotInDoubleQuotes) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return
	}
	a := f.(*action.Action)

	re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
	found := re.FindAllSubmatch(a.Raw, -1)
	for _, f := range found {
		chErrors <- fmt.Sprintf("action '%s' calls a variable '%s' that is in double quotes", a.DirName, string(f[1]))
		compliant = false
	}

	return
}
