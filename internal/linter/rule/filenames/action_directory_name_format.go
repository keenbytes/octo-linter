package filenames

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/casematch"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

// ActionDirectoryNameFormat checks if directory containing action adheres to the selected naming convention.
type ActionDirectoryNameFormat struct {
}

func (r ActionDirectoryNameFormat) ConfigName() string {
	return "filenames__action_directory_name_format"
}

func (r ActionDirectoryNameFormat) FileType() int {
	return rule.DotGithubFileTypeAction
}

func (r ActionDirectoryNameFormat) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != "dash-case" && val != "camelCase" && val != "PascalCase" && val != "ALL_CAPS" {
		return fmt.Errorf("value can be one of: dash-case, camelCase, PascalCase, ALL_CAPS")
	}

	return nil
}

func (r ActionDirectoryNameFormat) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	m := casematch.Match(a.DirName, conf.(string))
	if !m {
		chErrors <- fmt.Sprintf("action directory name '%s' must be %s", a.DirName, conf.(string))
		compliant = false
	}

	return
}
