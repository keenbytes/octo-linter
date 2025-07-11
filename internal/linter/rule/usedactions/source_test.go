package usedactions

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestSourceValidate(t *testing.T) {
	t.Parallel()

	rule := Source{}

	for _, confBad := range []interface{}{4, true, "wrong"} {
		err := rule.Validate(confBad)
		if err == nil {
			t.Errorf("Source.Validate should return error when conf is %v", confBad)
		}
	}

	for _, confGood := range []interface{}{"local-only", "local-or-external", "external-only"} {
		err := rule.Validate(confGood)
		if err != nil {
			t.Errorf("Source.Validate should not return error when conf is %v", confGood)
		}
	}
}

func TestLocalOnly(t *testing.T) {
	t.Parallel()

	rule := Source{}
	conf := "local-only"
	d := DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are external actions and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 3 {
			t.Errorf("Source.Lint on %s should send 3 errors over the channel not [%s]", n, strings.Join(ruleErrors, "\n"))
		}
	}

	ruletest.Action(d, "usedactions-source", fn)
	ruletest.Workflow(d, "usedactions-source.yml", fn)
}

func TestExternalOnlyOnAction(t *testing.T) {
	t.Parallel()

	rule := Source{}
	conf := "external-only"
	d := DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are local actions and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 2 {
			t.Errorf("Source.Lint on %s should send 2 errors over the channel not [%s]", n, strings.Join(ruleErrors, "\n"))
		}
	}

	ruletest.Action(d, "usedactions-source", fn)
	ruletest.Workflow(d, "usedactions-source.yml", fn)
}

func TestLocalOrExternalOnAction(t *testing.T) {
	t.Parallel()

	rule := Source{}
	conf := "local-or-external"
	d := DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are invalid actions that are nor local nor external, and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 1 {
			t.Errorf("Source.Lint on %s should send 1 error over the channel not [%s]", n, strings.Join(ruleErrors, "\n"))
		}
	}

	ruletest.Action(d, "usedactions-source", fn)
	ruletest.Workflow(d, "usedactions-source.yml", fn)
}
