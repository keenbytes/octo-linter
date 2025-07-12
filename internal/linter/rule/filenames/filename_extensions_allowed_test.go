package filenames

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestFilenameExtensionsAllowedValidate(t *testing.T) {
	t.Parallel()

	rule := FilenameExtensionsAllowed{}

	confBad := []interface{}{"something", "something2"}

	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("FilenameExtensionsAllowed.Validate should return error when conf is %v", confBad)
	}

	confGood := []interface{}{"yml", "yaml"}

	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf(
			"FilenameExtensionsAllowed.Validate should not return error (%s) when conf is %v",
			err.Error(),
			confGood,
		)
	}
}

func TestFilenameExtensionsAllowedNotCompliant(t *testing.T) {
	t.Parallel()

	rule := FilenameExtensionsAllowed{}
	conf := []interface{}{"yaml"}
	d := DotGithub

	fn := func(f dotgithub.File, _ string) {
		compliant, ruleErrors, err := ruletest.Lint(2, rule, conf, f, d)
		if compliant {
			t.Errorf(
				"FilenameExtensionsAllowed.Lint should return false when filename extension is not in config",
			)
		}

		if err != nil {
			t.Errorf("FilenameExtensionsAllowed.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) == 0 {
			t.Errorf("FilenameExtensionsAllowed.Lint should send an error over the channel")
		}
	}

	ruletest.Action(d, "valid-action", fn)
	ruletest.Workflow(d, "valid-workflow.yml", fn)
}

func TestFilenameExtensionsAllowedCompliant(t *testing.T) {
	t.Parallel()

	rule := FilenameExtensionsAllowed{}
	conf := []interface{}{"yml"}
	d := DotGithub

	fn := func(f dotgithub.File, _ string) {
		compliant, ruleErrors, err := ruletest.Lint(2, rule, conf, f, d)
		if !compliant {
			t.Errorf(
				"FilenameExtensionsAllowed.Lint should return true when filename extension is in config",
			)
		}

		if err != nil {
			t.Errorf("FilenameExtensionsAllowed.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) > 0 {
			t.Errorf(
				"FilenameExtensionsAllowed.Lint should not send any error over the channel, sent %s",
				strings.Join(ruleErrors, "|"),
			)
		}
	}

	ruletest.Action(d, "valid-action", fn)
	ruletest.Workflow(d, "valid-workflow.yml", fn)
}
