package filenames

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/casematch"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// ActionDirectoryNameFormat checks if directory containing action adheres to the selected naming convention.
type ActionDirectoryNameFormat struct {
}

func (r ActionDirectoryNameFormat) ConfigName(int) string {
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

func (r ActionDirectoryNameFormat) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	m := casematch.Match(a.DirName, conf.(string))
	if !m {
		chErrors <- glitch.Glitch{
			Path: a.Path,
			Name: a.DirName,
			Type: rule.DotGithubFileTypeAction,
			ErrText: fmt.Sprintf("directory name must be %s", conf.(string)),
			RuleName: r.ConfigName(0),
		}
		compliant = false
	}

	return
}
