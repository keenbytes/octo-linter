package rule

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// RuleWorkflowFileExtensions checks if workflow file extension is one of the specific values, eg. 'yml' or 'yaml'.
type RuleWorkflowFileExtensions struct {
	Value      []string
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowFileExtensions) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "yml" && v != "yaml" {
				return errors.New("workflow_file_extensions can only contain values of 'yml' and/or 'yaml'")
			}
		}
	}
	return nil
}

func (r RuleWorkflowFileExtensions) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	fileParts := strings.Split(w.FileName, ".")
	extension := fileParts[len(fileParts)-1]
	for _, v := range r.Value {
		if extension == v {
			return true, nil
		}
	}
	printErrOrWarn(r.ConfigName, r.IsError,
		fmt.Sprintf("workflow '%s' file extension must be one of: %s", w.DisplayName, strings.Join(r.Value, ",")),
		chWarnings, chErrors,
	)
	return false, nil
}

func (r RuleWorkflowFileExtensions) GetConfigName() string {
	return r.ConfigName
}
