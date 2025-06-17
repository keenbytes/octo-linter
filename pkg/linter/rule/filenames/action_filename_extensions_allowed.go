package rulefilenames

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/linter/rule"
)

// RuleActionFileExtensions checks if action file extension is one of the specific values, eg. 'yml' or 'yaml'.
type ActionFilenameExtensionsAllowed struct {
	rule.RuleAction
}

func (r ActionFilenameExtensionsAllowed) ConfigName() string {
	return "filenames__action_filename_extensions_allowed"
}

func (r ActionFilenameExtensionsAllowed) Validate(conf interface{}) error {
	vals, ok := conf.([]interface{})
	if !ok {
		return errors.New("value should be []string")
	}

	for _, v := range vals {
		extension, ok := v.(string)
		if !ok {
			return errors.New("value should be []string")
		}
		if extension != "yml" && extension != "yaml" {
			return fmt.Errorf("value can contain only 'yml' and/or 'yaml'")
		}
	}

	return nil
}

func (r ActionFilenameExtensionsAllowed) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	allowedExtensions, ok := conf.([]interface{})
	if !ok {
		return
	}

	pathParts := strings.Split(a.Path, "/")
	fileParts := strings.Split(pathParts[len(pathParts)-1], ".")
	extension := fileParts[len(fileParts)-1]

	var allowedExtensionsList []string
	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension.(string) {
			return
		}
		allowedExtensionsList = append(allowedExtensionsList, allowedExtension.(string))
	}
	compliant = false
	chErrors <- fmt.Sprintf("action '%s' file extension must be one of: %s", a.DirName, strings.Join(allowedExtensionsList, ","))

	return
}
