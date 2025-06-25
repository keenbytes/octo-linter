package ruletest

import (
	"errors"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func Lint(timeout int, rule rule.Rule, conf interface{}, f dotgithub.File, d *dotgithub.DotGithub) (compliant bool, err error, ruleErrors []string) {
	compliant = true

	timer := time.After(time.Duration(timeout) * time.Second)

	chErrors := make(chan string)

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
			case ruleError, more := <-chErrors:
				if more {
					ruleErrors = append(ruleErrors, ruleError)
				} else {
					break loop
				}
			}
		}

	return
}
