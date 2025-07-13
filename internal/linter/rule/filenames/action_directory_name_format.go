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
type ActionDirectoryNameFormat struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r ActionDirectoryNameFormat) ConfigName(int) string {
	return "filenames__action_directory_name_format"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r ActionDirectoryNameFormat) FileType() int {
	return rule.DotGithubFileTypeAction
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r ActionDirectoryNameFormat) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != ValueDashCase && val != ValueCamelCase && val != ValuePascalCase && val != ValueAllCaps {
		return fmt.Errorf("value can be one of: %s, %s, %s, %s", ValueDashCase, ValueCamelCase, ValuePascalCase, ValueAllCaps)
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r ActionDirectoryNameFormat) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if file.GetType() != rule.DotGithubFileTypeAction {
		return true, nil
	}

	actionInstance := file.(*action.Action)

	m := casematch.Match(actionInstance.DirName, conf.(string))
	if !m {
		chErrors <- glitch.Glitch{
			Path:     actionInstance.Path,
			Name:     actionInstance.DirName,
			Type:     rule.DotGithubFileTypeAction,
			ErrText:  "directory name must be " + conf.(string),
			RuleName: r.ConfigName(0),
		}

		return false, nil
	}

	return true, nil
}
