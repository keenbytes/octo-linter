package glitch

import (
	"fmt"
)

const (
	DotGithubFileTypeAction   = 1
	DotGithubFileTypeWorkflow = 2
)

type Glitch struct {
	Type     int
	Name     string
	Path     string
	RuleName string
	ErrText  string
	IsError  bool
}

func ListToMarkdown(glitches []*Glitch, limit int) (s string) {
	s = `|Item|Error|
|---|---|
`
	for i, g := range glitches {
		if limit > 0 && i == limit {
			break
		}
		name := fmt.Sprintf(`a/%s`, g.Name)
		if g.Type == DotGithubFileTypeWorkflow {
			name = fmt.Sprintf(`w/%s`, g.Name)
		}
		level := `🟠`
		if g.IsError {
			level = `🔴`
		}
		s += fmt.Sprintf("|%s|%s %s *(%s)*|\n", name, level, g.ErrText, g.RuleName)
	}

	if len(glitches) > limit && limit > 0 {
		s += fmt.Sprintf("\n...and many more (%d in total).", len(glitches))
	}

	return
}
