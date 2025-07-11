package workflow

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	DotGithubFileTypeWorkflow = 2
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

func (w *Workflow) Unmarshal(fromRaw bool) error {
	// TODO: fromRaw is not implemented
	pathSplit := strings.Split(w.Path, "/")
	w.FileName = pathSplit[len(pathSplit)-1]
	workflowName := strings.ReplaceAll(w.FileName, ".yaml", "")
	w.DisplayName = strings.ReplaceAll(workflowName, ".yml", "")

	slog.Debug(fmt.Sprintf("reading %s ...", w.Path))

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

func (w *Workflow) GetType() int {
	return DotGithubFileTypeWorkflow
}
