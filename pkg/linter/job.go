package linter

import (
	"fmt"
	"time"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/linter/rule"
)

type Job struct {
	rule rule.Rule
	file dotgithub.File
	dotGithub *dotgithub.DotGithub
	isError bool
}

func (j *Job) Run(chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true

	done := make(chan struct{})
	timer := time.NewTimer(time.Duration(10*time.Second))

	go func() {
		compliant, err = j.rule.Lint(j.file, j.dotGithub, chWarnings, chErrors)
		close(done)
	}()

	select {
	case <-timer.C:
		return false, fmt.Errorf("rule %s timed out", j.rule.GetConfigName())
	case <-done:
		return compliant, err
	}
}
