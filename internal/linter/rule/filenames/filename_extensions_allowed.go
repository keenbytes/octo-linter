package filenames

import (
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// FilenameExtensionsAllowed checks if file extension is one of the specific values, eg. 'yml' or 'yaml'.
type FilenameExtensionsAllowed struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
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

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r FilenameExtensionsAllowed) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r FilenameExtensionsAllowed) Validate(conf interface{}) error {
	vals, ok := conf.([]interface{})
	if !ok {
		return errValueNotStringArray
	}

	for _, v := range vals {
		extension, ok := v.(string)
		if !ok {
			return errValueNotStringArray
		}

		if extension != "yml" && extension != "yaml" {
			return errValueNotYmlOrYaml
		}
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r FilenameExtensionsAllowed) Lint(
	conf interface{},
	file dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	if file.GetType() != rule.DotGithubFileTypeAction &&
		file.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	allowedExtensions, confIsInterfaceArray := conf.([]interface{})
	if !confIsInterfaceArray {
		return false, errValueNotStringArray
	}

	var (
		extension    string
		filePath     string
		fileTypeName string
		fileType     int
	)

	if file.GetType() == rule.DotGithubFileTypeAction {
		actionInstance, ok := file.(*action.Action)
		if !ok {
			return false, errFileInvalidType
		}

		pathParts := strings.Split(actionInstance.Path, "/")
		fileParts := strings.Split(pathParts[len(pathParts)-1], ".")
		extension = fileParts[len(fileParts)-1]

		filePath = actionInstance.Path
		fileTypeName = actionInstance.DirName
		fileType = rule.DotGithubFileTypeAction
	}

	if file.GetType() == rule.DotGithubFileTypeWorkflow {
		workflowInstance, ok := file.(*workflow.Workflow)
		if !ok {
			return false, errFileInvalidType
		}

		fileParts := strings.Split(workflowInstance.FileName, ".")
		extension = fileParts[len(fileParts)-1]

		filePath = workflowInstance.Path
		fileTypeName = workflowInstance.DisplayName
		fileType = rule.DotGithubFileTypeWorkflow
	}

	allowedExtensionsList := make([]string, 0, len(allowedExtensions))

	for _, allowedExtensionInterface := range allowedExtensions {
		allowedExtension, ok := allowedExtensionInterface.(string)
		if !ok {
			return false, errValueNotStringArray
		}

		if extension == allowedExtension {
			return true, nil
		}

		allowedExtensionsList = append(allowedExtensionsList, allowedExtension)
	}

	chErrors <- glitch.Glitch{
		Path:     filePath,
		Name:     fileTypeName,
		Type:     fileType,
		ErrText:  "file extension must be one of: " + strings.Join(allowedExtensionsList, ","),
		RuleName: r.ConfigName(fileType),
	}

	return false, nil
}
