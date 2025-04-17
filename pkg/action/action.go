package action

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
)

const (
	DotGithubFileTypeAction = 1
)

type Action struct {
	Path        string
	Raw         []byte
	DirName     string
	Name        string                   `yaml:"name"`
	Description string                   `yaml:"description"`
	Inputs      map[string]*ActionInput  `yaml:"inputs"`
	Outputs     map[string]*ActionOutput `yaml:"outputs"`
	Runs        *ActionRuns              `yaml:"runs"`
}

func (a *Action) Unmarshal(logLevel int, fromRaw bool) error {
	if !fromRaw {
		if logLevel == loglevel.LogLevelDebug {
			fmt.Fprintf(os.Stdout, "dbg:reading %s ...\n", a.Path)
		}
		b, err := os.ReadFile(a.Path)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %w", a.Path, err)
		}
		a.Raw = b
	}
	err := yaml.Unmarshal(a.Raw, &a)
	if err != nil {
		return fmt.Errorf("cannot unmarshal file %s: %w", a.Path, err)
	}
	if a.Runs != nil {
		a.Runs.SetParentType("action")
	}
	return nil
}

func (a *Action) GetType() int {
	return DotGithubFileTypeAction
}
