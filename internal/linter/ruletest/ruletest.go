package ruletest

import (
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func RunLintAndGetRuleErrors(timeout int, rule rule.Rule, conf interface{}, f dotgithub.File, d *dotgithub.DotGithub) (compliant bool, err error, ruleErrors []string) {
	chErrors := make(chan string)
	timer := time.After(time.Duration(timeout) * time.Second)

	compliant = true
	err = nil

	go func() {
		compliant, err = rule.Lint(conf, f, d, chErrors)
	}()

	select {
	case <-timer:
		close(chErrors)
	case ruleError := <-chErrors:
		ruleErrors = append(ruleErrors, ruleError)
	}

	return compliant, err, ruleErrors
}
