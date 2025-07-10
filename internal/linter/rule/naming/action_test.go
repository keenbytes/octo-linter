package naming

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestActionValidate(t *testing.T) {
	rule := Action{}

	confBad := "some string"
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("Action.Validate should return error when conf is %v", confBad)
	}

	confGood := "camelCase"
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("Action.Validate should not return error (%s) when conf is %v", err.Error(), confGood)
	}
}

func TestActionNotCompliant(t *testing.T) {
	for field, conf := range map[string]string{
		"input_name": "dash-case",
		"output_name": "dash-case",
		"referenced_variable": "ALL_CAPS",
		"step_env": "ALL_CAPS",
	}{
		rule := Action{
			Field: field,
		}
		d := ruletest.DotGithub

		fn := func(f dotgithub.File, n string) {
			compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
			if compliant {
				t.Errorf("Action.Lint should return false when action %s does not follow naming convention of '%s'", field, conf)
			}
			if err != nil {
				t.Errorf("Action.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) != 2 {
				t.Errorf("Action.Lint should send 2 errors over the channel, got [%s]", strings.Join(ruleErrors, "\n"))
			}
		}

		ruletest.Action(d, "naming-action", fn)
	}
}

func TestActionCompliant(t *testing.T) {
	for field, conf := range map[string]string{
		"input_name": "dash-case",
		"output_name": "dash-case",
		"referenced_variable": "ALL_CAPS",
		"step_env": "ALL_CAPS",
	}{
		rule := Action{
			Field: field,
		}
		d := ruletest.DotGithub

		fn := func(f dotgithub.File, n string) {
			compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
			if !compliant {
				t.Errorf("Action.Lint should return true when action %s follows naming convention of '%s'", field, conf)
			}
			if err != nil {
				t.Errorf("Action.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) != 0 {
				t.Errorf("Action.Lint should not send any errors over the channel, got [%s]", strings.Join(ruleErrors, "\n"))
			}
		}

		ruletest.Action(d, "valid-action", fn)
	}
}
