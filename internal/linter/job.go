package linter

import (
	"fmt"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

type Job struct {
	rule      rule.Rule
	file      dotgithub.File
	dotGithub *dotgithub.DotGithub
	isError   bool
	value     interface{}
}

func (j *Job) Run(chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true

	done := make(chan struct{})
	timer := time.NewTimer(time.Duration(10 * time.Second))

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
		return false, fmt.Errorf("rule %s timed out", j.rule.ConfigName(j.file.GetType()))
	case <-done:
		return compliant, err
	}
}
