package workflow

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
)

type Workflow struct {
	Path        string
	Raw         []byte
	FileName    string
	DisplayName string
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description"`
	Env         map[string]string       `yaml:"env"`
	Jobs        map[string]*WorkflowJob `yaml:"jobs"`
	On          *WorkflowOn             `yaml:"on"`
}

func (w *Workflow) Unmarshal(logLevel int) error {
	pathSplit := strings.Split(w.Path, "/")
	w.FileName = pathSplit[len(pathSplit)-1]
	workflowName := strings.Replace(w.FileName, ".yaml", "", -1)
	w.DisplayName = strings.Replace(workflowName, ".yml", "", -1)

	if logLevel == loglevel.LogLevelDebug {
		fmt.Fprintf(os.Stdout, "dbg:reading %s ...\n", w.Path)
	}
	b, err := os.ReadFile(w.Path)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", w.Path, err)
	}
	w.Raw = b

	err = yaml.Unmarshal(w.Raw, &w)
	if err != nil {
		return fmt.Errorf("cannot unmarshal file %s: %w", w.Path, err)
	}
	if w.Jobs != nil {
		for _, j := range w.Jobs {
			j.SetParentType("workflow")
		}
	}
	return nil
}
