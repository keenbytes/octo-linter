package rule

import (
	"fmt"
	"strings"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionFileExtensions struct {
	Value      []string
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionFileExtensions) Validate() error {
	if len(r.Value) > 0 {
		for _, v := range r.Value {
			if v != "yml" && v != "yaml" {
				return fmt.Errorf("%s can only contain values of 'yml' and/or 'yaml'", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleActionFileExtensions) Lint(a *action.Action, d *dotgithub.DotGithub) (compliant bool, err error) {
	pathParts := strings.Split(a.Path, "/")
	fileParts := strings.Split(pathParts[len(pathParts)-1], ".")
	extension := fileParts[len(fileParts)-1]
	for _, v := range r.Value {
		if extension == v {
			return true, nil
		}
	}
	printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel,
		fmt.Sprintf("action '%s' file extension must be one of: %s", a.DirName, strings.Join(r.Value, ",")),
	)
	return false, nil
}

func (r RuleActionFileExtensions) GetConfigName() string {
	return r.ConfigName
}
