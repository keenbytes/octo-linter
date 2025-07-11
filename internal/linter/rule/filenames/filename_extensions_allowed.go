package filenames

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// FilenameExtensionsAllowed checks if file extension is one of the specific values, eg. 'yml' or 'yaml'.
type FilenameExtensionsAllowed struct {
}

func (r FilenameExtensionsAllowed) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "filenames__workflow_filename_extensions_allowed"
	case rule.DotGithubFileTypeAction:
		return "filenames__action_filename_extensions_allowed"
	default:
		return "filenames__*_filename_extensions_allowed*__not_in_double_quotes"
	}
}

func (r FilenameExtensionsAllowed) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r FilenameExtensionsAllowed) Validate(conf interface{}) error {
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

func (r FilenameExtensionsAllowed) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	allowedExtensions, ok := conf.([]interface{})
	if !ok {
		return true, nil
	}

	var (
		extension    string
		filePath     string
		fileTypeName string
		fileType     int
	)

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)

		pathParts := strings.Split(a.Path, "/")
		fileParts := strings.Split(pathParts[len(pathParts)-1], ".")
		extension = fileParts[len(fileParts)-1]

		filePath = a.Path
		fileTypeName = a.DirName
		fileType = rule.DotGithubFileTypeAction
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)

		fileParts := strings.Split(w.FileName, ".")
		extension = fileParts[len(fileParts)-1]

		filePath = w.Path
		fileTypeName = w.DisplayName
		fileType = rule.DotGithubFileTypeWorkflow
	}

	var allowedExtensionsList []string

	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension.(string) {
			return true, nil
		}

		allowedExtensionsList = append(allowedExtensionsList, allowedExtension.(string))
	}

	chErrors <- glitch.Glitch{
		Path:     filePath,
		Name:     fileTypeName,
		Type:     fileType,
		ErrText:  fmt.Sprintf("file extension must be one of: %s", strings.Join(allowedExtensionsList, ",")),
		RuleName: r.ConfigName(fileType),
	}

	return false, nil
}
