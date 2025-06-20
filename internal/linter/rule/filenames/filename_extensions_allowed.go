package filenames

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

func (r FilenameExtensionsAllowed) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	if f.GetType() == rule.DotGithubFileTypeAction {
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
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)

		re := regexp.MustCompile(`\"\${{[ ]*([a-zA-Z0-9\\-_.]+)[ ]*}}\"`)
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			chErrors <- fmt.Sprintf("workflow '%s' calls a variable '%s' that is in double quotes", w.FileName, string(f[1]))
			compliant = false
		}
	}

	return
}
