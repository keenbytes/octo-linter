package filenames

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowFilenameExtensionsAllowed checks if workflow file extension is one of the specific values, eg. 'yml' or 'yaml'.
type WorkflowFilenameExtensionsAllowed struct {
}

func (r WorkflowFilenameExtensionsAllowed) ConfigName() string {
	return "filenames__workflow_filename_extensions_allowed"
}

func (r WorkflowFilenameExtensionsAllowed) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowFilenameExtensionsAllowed) Validate(conf interface{}) error {
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

func (r WorkflowFilenameExtensionsAllowed) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	w := f.(*workflow.Workflow)

	allowedExtensions, ok := conf.([]interface{})
	if !ok {
		return
	}

	fileParts := strings.Split(w.FileName, ".")
	extension := fileParts[len(fileParts)-1]

	var allowedExtensionsList []string
	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension.(string) {
			return
		}
		allowedExtensionsList = append(allowedExtensionsList, allowedExtension.(string))
	}
	compliant = false
	chErrors <- fmt.Sprintf("workflow '%s' file extension must be one of: %s", w.DisplayName, strings.Join(allowedExtensionsList, ","))

	return
}
