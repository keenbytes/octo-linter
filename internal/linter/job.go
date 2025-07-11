package linter

import (
	"fmt"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

func errRuleTimeout(name string) error {
	return fmt.Errorf("rule %s timed out", name)
}

type Job struct {
	rule      rule.Rule
	file      dotgithub.File
	dotGithub *dotgithub.DotGithub
	isError   bool
	value     interface{}
}

func (j *Job) Run(chWarnings chan<- glitch.Glitch, chErrors chan<- glitch.Glitch) (bool, error) {
	compliant := true

	var err error

	done := make(chan struct{})
	timer := time.NewTimer(10 * time.Second)

	go func() {
		if j.isError {
			compliant, err = j.rule.Lint(j.value, j.file, j.dotGithub, chErrors)
		} else {
			compliant, err = j.rule.Lint(j.value, j.file, j.dotGithub, chWarnings)
		}

		close(done)
	}()

	select {
	case <-timer.C:
		return false, errRuleTimeout(j.rule.ConfigName(j.file.GetType()))
	case <-done:
		return compliant, fmt.Errorf("rule lint error: %w", err)
	}
}
