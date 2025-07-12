// Package ruletest contains helper functions for testing rules.
package ruletest

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// Lint runs a rule with specific configuration on a specified file and returns all lint errors and a boolean indicating whether it is compliant or not.
func Lint(
	timeout int,
	rule rule.Rule,
	conf interface{},
	file dotgithub.File,
	dotGithub *dotgithub.DotGithub,
) (bool, []string, error) {
	compliant := true
	ruleErrors := []string{}

	var err error

	timer := time.After(time.Duration(timeout) * time.Second)

	chErrors := make(chan glitch.Glitch)

	go func() {
		compliant, err = rule.Lint(conf, file, dotGithub, chErrors)
		close(chErrors)
	}()

loop:
	for {
		select {
		case <-timer:
			err = errors.New("timeout")
			compliant = false

			break loop
		case glitchInstance, more := <-chErrors:
			if more {
				ruleError := fmt.Sprintf("%s %s: %s", glitchInstance.Path, glitchInstance.RuleName, glitchInstance.ErrText)
				ruleErrors = append(ruleErrors, ruleError)
			} else {
				break loop
			}
		}
	}

	return compliant, ruleErrors, err
}

// Action runs a test function on a specific action in DotGithub.
func Action(
	dotGithub *dotgithub.DotGithub,
	actionToTest string,
	testFunc func(file dotgithub.File, name string),
) {
	for actionName, actionFile := range dotGithub.Actions {
		if actionName != actionToTest {
			continue
		}

		log.Printf("running test on action %s...", actionName)
		testFunc(actionFile, actionName)
	}
}

// Workflow runs a test function on a specific workflow from DotGithub.
func Workflow(
	dotGithub *dotgithub.DotGithub,
	workflowToTest string,
	testFunc func(file dotgithub.File, name string),
) {
	for workflowName, workflowFile := range dotGithub.Workflows {
		if workflowName != workflowToTest {
			continue
		}

		log.Printf("running test on workflow %s...", workflowName)
		testFunc(workflowFile, workflowName)
	}
}
