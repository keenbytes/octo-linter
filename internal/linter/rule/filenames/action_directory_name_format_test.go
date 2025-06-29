package filenames

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestActionDirectoryNameFormatValidate(t *testing.T) {
	rule := ActionDirectoryNameFormat{}

	confBad := "some string"
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("ActionDirectoryNameFormat.Validate should return error when conf is %v", confBad)
	}

	confGood := "camelCase"
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("ActionDirectoryNameFormat.Validate should not return error (%s) when conf is %v", err.Error(), confGood)
	}
}

func TestActionDirectoryNameFormatNotCompliant(t *testing.T) {
	rule := ActionDirectoryNameFormat{}
	d := ruletest.DotGithub

	for _, nameFormat := range []string{"camelCase", "PascalCase", "ALL_CAPS"}{
		fn := func(f dotgithub.File, n string) {
			compliant, err, ruleErrors := ruletest.Lint(2, rule, nameFormat, f, d)
			if compliant {
				t.Errorf("ActionDirectoryNameFormat.Lint should return false when filename is not %s", nameFormat)
			}
			if err != nil {
				t.Errorf("ActionDirectoryNameFormat.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) == 0 {
				t.Errorf("ActionDirectoryNameFormat.Lint should send an error over the channel when filename is not %s", nameFormat)
			}
		}

		ruletest.Action(d, "valid-action", fn)
	}
}

func TestActionDirectoryNameFormatCompliant(t *testing.T) {
	rule := ActionDirectoryNameFormat{}
	conf := "dash-case"
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
		if !compliant {
			t.Errorf("ActionDirectoryNameFormat.Lint should return true when filename is %s", conf)
		}
		if err != nil {
			t.Errorf("ActionDirectoryNameFormat.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) > 0 {
			t.Errorf("ActionDirectoryNameFormat.Lint should not send any error over the channel, sent %s", strings.Join(ruleErrors, "|"))
		}
	}

	ruletest.Action(d, "valid-action", fn)
}
