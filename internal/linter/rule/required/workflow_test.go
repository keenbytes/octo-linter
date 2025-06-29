package required

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func TestWorkflowValidate(t *testing.T) {
	rule := Workflow{
		Field: "workflow",
	}

	confBad := 4
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("Workflow.Validate should return error when conf is not []string")
	}

	confGood := []interface{}{"name"}
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("Workflow.Validate should not return error when conf is []string")
	}

	for _, f := range []string{"dispatch_input", "call_input"}{
		rule = Workflow{
			Field: f,
		}

		confBad2 := []interface{}{"name", "description"}
		err = rule.Validate(confBad2)
		if err == nil {
			t.Errorf("Workflow.Validate should return error when conf contains invalid values")
		}

		confGood2 := []interface{}{"description"}
		err = rule.Validate(confGood2)
		if err != nil {
			t.Errorf("Workflow.Validate should not return error when conf contains valid values")
		}
	}
}
func TestWorkflowFieldWorkflowNotCompliant(t *testing.T) {
	rule := Workflow{
		Field: "workflow",
	}
	conf := []interface{}{"name"}
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
		if compliant {
			t.Errorf("Workflow.Lint should return false when workflow does not have a 'name' field")
		}
		if err != nil {
			t.Errorf("Workflow.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) == 0 {
			t.Errorf("Workflow.Lint should send an error over the channel")
		}
	}

	ruletest.Workflow(d, "required-workflow.yml", fn)
}

func TestWorkflowFieldCallInputDispatchInputNotCompliant(t *testing.T) {
	for _, field := range []string{"dispatch_input", "call_input"}{
		rule := Workflow{
			Field: field,
		}
		conf := []interface{}{"description"}
		d := ruletest.DotGithub

		fn := func(f dotgithub.File, n string) {
			compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
			if compliant {
				t.Errorf("Workflow.Lint should return false when workflow %s does not have a 'description' field", field)
			}
			if err != nil {
				t.Errorf("Workflow.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) != 2 {
				t.Errorf("Workflow.Lint should send 2 errors over the channel, got [%s]", strings.Join(ruleErrors, "\n"))
			}
		}

		ruletest.Workflow(d, "required-workflow.yml", fn)
	}
}

func TestWorkflowFieldWorkflowCompliant(t *testing.T) {
	rule := Workflow{
		Field: "workflow",
	}
	conf := []interface{}{"name"}
	d := ruletest.DotGithub

	fn := func(f dotgithub.File, n string) {
		compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
		if !compliant {
			t.Errorf("Workflow.Lint should return true when workflow has a 'name' field")
		}
		if err != nil {
			t.Errorf("Workflow.Lint failed with an error: %s", err.Error())
		}

		if len(ruleErrors) != 0 {
			t.Errorf("Workflow.Lint should not send any error over the channel, got [%s]", strings.Join(ruleErrors, "\n"))
		}
	}

	ruletest.Workflow(d, "valid-workflow.yml", fn)
}

func TestWorkflowFieldCallInputDispatchInputCompliant(t *testing.T) {
	for _, field := range []string{"dispatch_input", "call_input"}{
		rule := Workflow{
			Field: field,
		}
		conf := []interface{}{"description"}
		d := ruletest.DotGithub

		fn := func(f dotgithub.File, n string) {
			compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
			if !compliant {
				t.Errorf("Workflow.Lint should return true when workflow %s has a 'description' field", field)
			}
			if err != nil {
				t.Errorf("Workflow.Lint failed with an error: %s", err.Error())
			}

			if len(ruleErrors) != 0 {
				t.Errorf("Workflow.Lint should not send any errors over the channel, got [%s]", strings.Join(ruleErrors, "\n"))
			}
		}

		ruletest.Workflow(d, "valid-workflow.yml", fn)
	}
}
