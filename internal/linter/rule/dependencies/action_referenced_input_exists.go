package dependencies

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// ActionReferencedInputExists scans the action code for all input references and verifies that each has been previously defined.
// During action execution, if a reference to an undefined input is found, it is replaced with an empty string.
type ActionReferencedInputExists struct {
}

func (r ActionReferencedInputExists) ConfigName() string {
	return "dependencies__action_referenced_input_must_exists"
}

func (r ActionReferencedInputExists) FileType() int {
	return rule.DotGithubFileTypeAction
}

func (r ActionReferencedInputExists) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r ActionReferencedInputExists) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return
	}
	a := f.(*action.Action)

	re := regexp.MustCompile(`\${{[ ]*inputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
	found := re.FindAllSubmatch(a.Raw, -1)
	for _, f := range found {
		if a.Inputs == nil || a.Inputs[string(f[1])] == nil {
			chErrors <- fmt.Sprintf("action '%s' calls an input '%s' that does not exist", a.DirName, string(f[1]))
			compliant = false
		}
	}

	return
}
