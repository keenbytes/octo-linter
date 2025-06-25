package runners

import (
	"strings"
	"testing"

	"github.com/keenbytes/octo-linter/v2/internal/linter/ruletest"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

func TestNotLatestValidate(t *testing.T) {
	rule := NotLatest{}

	confBad := 4
	err := rule.Validate(confBad)
	if err == nil {
		t.Errorf("NotLatest.Validate should return error when conf is not bool")
	}

	confGood := true
	err = rule.Validate(confGood)
	if err != nil {
		t.Errorf("NotLatest.Validate should not return error when conf is bool")
	}
}

func TestNotLatestNotCompliant(t *testing.T) {
	rule := NotLatest{}
	conf := true
	d := &dotgithub.DotGithub{}
	f := &workflow.Workflow{
		FileName: "workflow.yml",
		Jobs: map[string]*workflow.WorkflowJob{
			"job1": {
				RunsOn: "ubuntu-22.04",
			},
			"job2": {
				RunsOn: "ubuntu-latest",
			},
		},
	}

	compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
	if compliant {
		t.Errorf("NotLatest.Lint should return false when 'latest' is found in at least one job")
	}
	if err != nil {
		t.Errorf("NotLatest.Lint failed with an error: %s", err.Error())
	}

	if len(ruleErrors) == 0 {
		t.Errorf("NotLatest.Lint should send an error over the channel")
	}
}

func TestNotLatestCompliant(t *testing.T) {
	rule := NotLatest{}
	conf := true
	d := &dotgithub.DotGithub{}
	f := &workflow.Workflow{
		FileName: "workflow.yml",
		Jobs: map[string]*workflow.WorkflowJob{
			"job1": {
				RunsOn: "ubuntu-22.04",
			},
			"job2": {
				RunsOn: "ubuntu-24.04",
			},
		},
	}

	compliant, err, ruleErrors := ruletest.Lint(2, rule, conf, f, d)
	if !compliant {
		t.Errorf("NotLatest.Lint should return true when 'latest' is not in any job")
	}
	if err != nil {
		t.Errorf("NotLatest.Lint failed with an error: %s", err.Error())
	}

	if len(ruleErrors) > 0 {
		t.Errorf("NotLatest.Lint should not send any error over the channel, sent: %s", strings.Join(ruleErrors, ","))
	}
}
