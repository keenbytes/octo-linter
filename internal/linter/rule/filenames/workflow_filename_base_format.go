package filenames

import (
	"errors"
	"fmt"
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/casematch"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// WorkflowFilenameBaseFormat checks if workflow file basename (without extension) adheres to the selected naming convention.
type WorkflowFilenameBaseFormat struct {
}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r WorkflowFilenameBaseFormat) ConfigName(int) string {
	return "filenames__workflow_filename_base_format"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r WorkflowFilenameBaseFormat) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

// Validate checks whether the given value is valid for this rule's configuration.
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

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r WorkflowFilenameBaseFormat) Lint(conf interface{}, f dotgithub.File, _ *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return true, nil
	}

	w := f.(*workflow.Workflow)

	fileParts := strings.Split(w.FileName, ".")
	basename := fileParts[0]

	m := casematch.Match(basename, conf.(string))
	if !m {
		chErrors <- glitch.Glitch{
			Path:     w.Path,
			Name:     w.DisplayName,
			Type:     rule.DotGithubFileTypeWorkflow,
			ErrText:  fmt.Sprintf("filename base must be %s", conf.(string)),
			RuleName: r.ConfigName(0),
		}

		return false, nil
	}

	return true, nil
}
