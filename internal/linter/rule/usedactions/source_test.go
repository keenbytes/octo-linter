package usedactions

import (
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/step"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

var TestStepsLocalOnly = []*step.Step{
	{
		Uses: "./.github/actions/validAction",
	},
	{
		Uses: "external-action/that-is-not-allowed",
	},
	{
		Uses: "another-external/not-allowed@v2",
	},
	{
		Uses: "./.github/actions/validActionAgain",
	},
}

var TestStepsExternalOnly = []*step.Step{
	{
		Uses: "./.github/actions/validActionNotAllowed",
	},
	{
		Uses: "external-org/repo/action-allowed-1@v1.0.0",
	},
	{
		Uses: "external-org/repo/action-allowed-2@v2",
	},
}

var TestStepsLocalOrExternal = []*step.Step{
	{
		Uses: "./.github/actions/validActionNotAllowed",
	},
	{
		Uses: "external-org/repo/action-allowed-1@v1.0.0",
	},
	{
		Uses: "./.github/some-wrong-name",
	},
	{
		Uses: "org/repo-that-is-wrong",
	},
	{
		Uses: "some-wrong-name-2",
	},
}

func TestSourceValidate(t *testing.T) {
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
	rule := Source{}
	conf := "local-only"
	d := &dotgithub.DotGithub{}

	for n, f := range map[string]dotgithub.File{
		"action": &action.Action{
			DirName: "action1",
			Runs: &action.ActionRuns{
				Steps: TestStepsLocalOnly,
			},
		},
		"workflow": &workflow.Workflow{
			FileName: "workflow1.yml",
			Jobs: map[string]*workflow.WorkflowJob{
				"main": {
					Steps: TestStepsLocalOnly,
				},
			},
		},
	} {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are external actions and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 2 {
			t.Errorf("Source.Lint on %s should send 2 errors over the channel not %v", n, ruleErrors)
		}
	}
}

func TestExternalOnlyOnAction(t *testing.T) {
	rule := Source{}
	conf := "external-only"
	d := &dotgithub.DotGithub{}
	for n, f := range map[string]dotgithub.File{
		"action": &action.Action{
			DirName: "action1",
			Runs: &action.ActionRuns{
				Steps: TestStepsExternalOnly,
			},
		},
		"workflow": &workflow.Workflow{
			FileName: "workflow1.yml",
			Jobs: map[string]*workflow.WorkflowJob{
				"main": {
					Steps: TestStepsExternalOnly,
				},
			},
		},
	} {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are local actions and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 1 {
			t.Errorf("Source.Lint on %s should send 2 errors over the channel not %v", n, ruleErrors)
		}
	}
}

func TestLocalOrExternalOnAction(t *testing.T) {
	rule := Source{}
	conf := "local-or-external"
	d := &dotgithub.DotGithub{}
	for n, f := range map[string]dotgithub.File{
		"action": &action.Action{
			DirName: "action1",
			Runs: &action.ActionRuns{
				Steps: TestStepsLocalOrExternal,
			},
		},
		"workflow": &workflow.Workflow{
			FileName: "workflow1.yml",
			Jobs: map[string]*workflow.WorkflowJob{
				"main": {
					Steps: TestStepsLocalOrExternal,
				},
			},
		},
	} {
		compliant, err, ruleErrors := ruletest.Lint(3, rule, conf, f, d)
		if compliant {
			t.Errorf("Source.Lint on %s should return false when there are invalid actions that are nor local nor external, and conf is %v", n, conf)
		}
		if err != nil {
			t.Errorf("Source.Lint on %s failed with an error: %s", n, err.Error())
		}

		if len(ruleErrors) != 3 {
			t.Errorf("Source.Lint on %s should send 2 errors over the channel not %v", n, ruleErrors)
		}
	}
}
