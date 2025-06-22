package usedactions

import (
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/step"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

func TestLocalOnlyOnWorkflow(t *testing.T) {
	rule := Source{}
	conf := "local-only"
	d := &dotgithub.DotGithub{}
	f := &workflow.Workflow{
		FileName: "workflow1.yml",
		Jobs: map[string]*workflow.WorkflowJob{
			"main": {
				Steps: []*step.Step{
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
				},
			},
		},
	}

	compliant, err, ruleErrors := ruletest.RunLintAndGetRuleErrors(3, rule, conf, f, d)
	if compliant {
		t.Errorf("Source.Lint should return false when there are external actions and conf is %v", conf)
	}
	if err != nil {
		t.Errorf("Source.Lint failed with an error: %s", err.Error())
	}

	if len(ruleErrors) != 2 {
		t.Errorf("Source.Lint should send 2 errors over the channel not %v", ruleErrors)
	}
}

func TestExternalOnlyOnWorkflow(t *testing.T) {
	rule := Source{}
	conf := "external-only"
	d := &dotgithub.DotGithub{}
	f := &workflow.Workflow{
		FileName: "workflow1.yml",
		Jobs: map[string]*workflow.WorkflowJob{
			"main": {
				Steps: []*step.Step{
					{
						Uses: "./.github/actions/validActionNotAllowed",
					},
					{
						Uses: "external-org/repo/action-allowed-1@v1.0.0",
					},
					{
						Uses: "external-org/repo/action-allowed-2@v2",
					},
				},
			},
		},
	}

	compliant, err, ruleErrors := ruletest.RunLintAndGetRuleErrors(3, rule, conf, f, d)
	if compliant {
		t.Errorf("Source.Lint should return false when there are local actions and conf is %v", conf)
	}
	if err != nil {
		t.Errorf("Source.Lint failed with an error: %s", err.Error())
	}

	if len(ruleErrors) != 1 {
		t.Errorf("Source.Lint should send 2 errors over the channel not %v", ruleErrors)
	}
}

func TestLocalOrExternalOnWorkflow(t *testing.T) {
	rule := Source{}
	conf := "local-or-external"
	d := &dotgithub.DotGithub{}

	f := &workflow.Workflow{
		FileName: "workflow1.yml",
		Jobs: map[string]*workflow.WorkflowJob{
			"main": {
				Steps: []*step.Step{
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
				},
			},
		},
	}

	compliant, err, ruleErrors := ruletest.RunLintAndGetRuleErrors(3, rule, conf, f, d)
	if compliant {
		t.Errorf("Source.Lint should return false when there are invalid actions that are nor local nor external, and conf is %v", conf)
	}
	if err != nil {
		t.Errorf("Source.Lint failed with an error: %s", err.Error())
	}

	if len(ruleErrors) != 3 {
		t.Errorf("Source.Lint should send 2 errors over the channel not %v", ruleErrors)
	}
}
