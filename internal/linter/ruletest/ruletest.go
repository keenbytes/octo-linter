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

func Lint(timeout int, rule rule.Rule, conf interface{}, f dotgithub.File, d *dotgithub.DotGithub) (compliant bool, err error, ruleErrors []string) {
	compliant = true

	timer := time.After(time.Duration(timeout) * time.Second)

	chErrors := make(chan glitch.Glitch)

	go func() {
		compliant, err = rule.Lint(conf, f, d, chErrors)
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

	return
}

func Action(d *dotgithub.DotGithub, action string, fn func(f dotgithub.File, n string)) {
	for n, f := range d.Actions {
		if n != action {
			continue
		}
		log.Printf("running test on action %s...", n)
		fn(f, n)
	}
}

func Workflow(d *dotgithub.DotGithub, workflow string, fn func(f dotgithub.File, n string)) {
	for n, f := range d.Workflows {
		if n != workflow {
			continue
		}
		log.Printf("running test on workflow %s...", n)
		fn(f, n)
	}
}
