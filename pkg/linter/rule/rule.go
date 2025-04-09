package rule

import (
	"fmt"
	"os"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type ActionRule interface {
	Validate() error
	Lint(a *action.Action, d *dotgithub.DotGithub) (bool, error)
	GetConfigName() string
}

type WorkflowRule interface {
	Validate() error
	Lint(a *workflow.Workflow, d *dotgithub.DotGithub) (bool, error)
	GetConfigName() string
}

type DotGithubRule interface {
	Validate() error
	Lint(d *dotgithub.DotGithub) (bool, error)
	GetConfigName() string
}

func printErrOrWarn(configName string, isError bool, logLevel int, errStr string) {
	if logLevel == loglevel.LogLevelNone {
		return
	}
	if isError && logLevel != loglevel.LogLevelNone {
		fmt.Fprintf(os.Stderr, "err:%s: %s\n", configName, errStr)
		return
	}
	if !isError && logLevel != loglevel.LogLevelOnlyErrors {
		fmt.Fprintf(os.Stderr, "wrn:%s: %s\n", configName, errStr)
		return
	}
}
