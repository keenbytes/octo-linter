package usedactions

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestExistsValidate(t *testing.T) {
	rule := Exists{}

	confBad := []interface{}{"something", "something2"}
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("Exists.Validate should return error when conf is %v", confBad)
	}

	confGood := []interface{}{"local", "external"}
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("Exists.Validate should not return error (%s) when conf is %v", err.Error(), confGood)
	}
}

func TestLocal(t *testing.T) {
	rule := Exists{}
	conf := []interface{}{"local"}
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Exists.Lint on %s should return false when conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Exists.Lint on %s failed with an error: %s", n, err.Error())
		}
		if len(ruleErrors) != 2 {
			t.Errorf("Exists.Lint on %s should send 2 errors over the channel not %s", n, strings.Join(ruleErrors, "|"))
		}
	}

	ruletest.Action(d, "usedactions-exists-local", fn)
	ruletest.Workflow(d, "usedactions-exists-local.yml", fn)
}

func TestExternal(t *testing.T) {
	rule := Exists{}
	conf := []interface{}{"external"}
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Exists.Lint on %s should return false when conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Exists.Lint on %s failed with an error: %s", n, err.Error())
		}
		if len(ruleErrors) != 2 {
			t.Errorf("Exists.Lint on %s should send 2 errors over the channel not %s", n, strings.Join(ruleErrors, "|"))
		}
	}

	ruletest.Action(d, "usedactions-exists-external", fn)
	ruletest.Workflow(d, "usedactions-exists-external.yml", fn)
}
