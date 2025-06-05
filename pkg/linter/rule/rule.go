package rule

import (
	"fmt"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
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

func printErrOrWarn(configName string, isError bool, errStr string, chWarnings chan<- string, chErrors chan<- string) {
	if isError {
		chErrors <- fmt.Sprintf("%s: %s", configName, errStr)
		return
	}
	if !isError {
		chWarnings <- fmt.Sprintf("%s: %s", configName, errStr)
		return
	}
}
