package usedactions

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestValidInputsValidate(t *testing.T) {
	rule := ValidInputs{}

	confBad := 4
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("ValidInputs.Validate should return error when conf is %v", confBad)
	}

	confGood := true
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("ValidInputs.Validate should not return error (%s) when conf is %v", err.Error(), confGood)
	}
}

func TestValidInputsNotCompliant(t *testing.T) {
	rule := ValidInputs{}
	conf := true
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
		if compliant {
			t.Errorf("ValidInputs.Lint should return false when there invalid inputs used when calling an action")
		}
		if err != nil {
			t.Errorf("ValidInputs.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) == 0 {
			t.Errorf("ValidInputs.Lint should send an error over the channel")
		}
	}

	ruletest.Workflow(d, "usedactions-valid-inputs.yml", fn)
}

func TestValidInputsCompliant(t *testing.T) {
	rule := ValidInputs{}
	conf := true
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
		if !compliant {
			t.Errorf("ValidInputs.Lint should return true when there are not any invalid inputs")
		}
		if err != nil {
			t.Errorf("ValidInputs.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) > 0 {
			t.Errorf("ValidInputs.Lint should not send any error over the channel, sent %s", strings.Join(ruleErrors, "|"))
		}
	}

	ruletest.Workflow(d, "valid-workflow.yml", fn)
}
