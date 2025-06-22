package runners

import (
	"testing"
	"time"

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

	chErrors := make(chan string)
	ruleError := ""
	timeout := time.After(2 * time.Second)

	go func() {
		compliant, err := rule.Lint(true, f, d, chErrors)
		if compliant {
			t.Errorf("NotLatest.Lint should return false when 'latest' is found in at least one job")
		}
		if err != nil {
			t.Errorf("NotLatest.Lint failed with an error")
		}
	}()

	select {
	case <-timeout:
		close(chErrors)
	case ruleError = <-chErrors:
	}

	if ruleError == "" {
		t.Errorf("NotLatest.Lint should send an error over the channel")
	}
}

func TestNotLatestCompliant(t *testing.T) {
	rule := NotLatest{}
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

	chErrors := make(chan string)
	ruleError := ""
	timeout := time.After(2 * time.Second)

	go func() {
		compliant, err := rule.Lint(true, f, d, chErrors)
		if !compliant {
			t.Errorf("NotLatest.Lint should return true when 'latest' is not found in any job")
		}
		if err != nil {
			t.Errorf("NotLatest.Lint failed with an error")
		}
	}()

	select {
	case <-timeout:
		close(chErrors)
	case ruleError = <-chErrors:
	}

	if ruleError != "" {
		t.Errorf("NotLatest.Lint should not send any error over the channel, sent %s", ruleError)
	}
}
