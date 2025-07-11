package filenames

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestWorkflowFilenameBaseFormatValidate(t *testing.T) {
	t.Parallel()

	rule := WorkflowFilenameBaseFormat{}

	confBad := "some string"

	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("WorkflowFilenameBaseFormat.Validate should return error when conf is %v", confBad)
	}

	confGood := "camelCase"

	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("WorkflowFilenameBaseFormat.Validate should not return error (%s) when conf is %v", err.Error(), confGood)
	}
}

func TestWorkflowFilenameBaseFormatNotCompliant(t *testing.T) {
	t.Parallel()

	rule := WorkflowFilenameBaseFormat{}
	d := DotGithub

	for _, nameFormat := range []string{"camelCase", "PascalCase", "ALL_CAPS"} {
		fn := func(f dotgithub.File, n string) {
			compliant, ruleErrors, err := ruletest.Lint(2, rule, nameFormat, f, d)
			if compliant {
				t.Errorf("WorkflowFilenameBaseFormat.Lint should return false when filename is not %s", nameFormat)
			}

			if err != nil {
				t.Errorf("WorkflowFilenameBaseFormat.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) == 0 {
				t.Errorf("WorkflowFilenameBaseFormat.Lint should send an error over the channel when filename is not %s", nameFormat)
			}
		}

		ruletest.Workflow(d, "valid-workflow.yml", fn)
	}
}

func TestWorkflowFilenameBaseFormatCompliant(t *testing.T) {
	t.Parallel()

	rule := WorkflowFilenameBaseFormat{}
	conf := "dash-case"
	d := DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, ruleErrors, err := ruletest.Lint(2, rule, conf, f, d)
		if !compliant {
			t.Errorf("WorkflowFilenameBaseFormat.Lint should return true when filename is %s", conf)
		}

		if err != nil {
			t.Errorf("WorkflowFilenameBaseFormat.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) > 0 {
			t.Errorf("WorkflowFilenameBaseFormat.Lint should not send any error over the channel, sent %s", strings.Join(ruleErrors, "|"))
		}
	}

	ruletest.Workflow(d, "valid-workflow.yml", fn)
}
