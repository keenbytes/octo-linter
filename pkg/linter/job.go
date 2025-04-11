package linter

import (
	"context"
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/linter/rule"
)

type Job struct {
	rule rule.Rule
	file dotgithub.File
	dotGithub *dotgithub.DotGithub
	isError bool
}

func (j *Job) Run(ctx context.Context, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true

	done := make(chan struct{})

	go func() {
		compliant, err = j.rule.Lint(j.file, j.dotGithub, chWarnings, chErrors)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("rule %s timed out", j.rule.GetConfigName())
	case <-done:
		return compliant, err
	}
}
