package filenames

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/casematch"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowFilenameBaseFormat checks if workflow file basename (without extension) adheres to the selected naming convention.
type WorkflowFilenameBaseFormat struct {
}

func (r WorkflowFilenameBaseFormat) ConfigName() string {
	return "filenames__workflow_filename_base_format"
}

func (r WorkflowFilenameBaseFormat) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowFilenameBaseFormat) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != "dash-case" && val != "dash-case;underscore-prefix-allowed" && val != "camelCase" && val != "PascalCase" && val != "ALL_CAPS" {
		return fmt.Errorf("value can be one of: dash-case, dash-case;underscore-prefix-allowed, camelCase, PascalCase, ALL_CAPS")
	}

	return nil
}

func (r WorkflowFilenameBaseFormat) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	fileParts := strings.Split(w.FileName, ".")
	basename := fileParts[0]

	m := casematch.Match(basename, conf.(string))
	if !m {
		chErrors <- fmt.Sprintf("workflow filename base '%s' must be %s", basename, conf.(string))
		compliant = false
	}

	return
}
