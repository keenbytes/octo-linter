package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
)

const (
	DotGithubFileTypeAction   = 1
	DotGithubFileTypeWorkflow = 2
)

type Rule interface {
	Validate() error
	Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (bool, error)
	GetConfigName() string
}

func printErrOrWarn(configName string, isError bool, logLevel int, errStr string, chWarnings chan<- string, chErrors chan<- string) {
	if logLevel == loglevel.LogLevelNone {
		return
	}
	if isError && logLevel != loglevel.LogLevelNone {
		chErrors <- fmt.Sprintf("%s: %s", configName, errStr)
		return
	}
	if !isError && logLevel != loglevel.LogLevelOnlyErrors {
		chWarnings <- fmt.Sprintf("%s: %s", configName, errStr)
		return
	}
}
