// Package glitch contains code related to representing a lint error.
package glitch

import (
	"fmt"
)

const (
	// DotGithubFileTypeAction represents the action file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeAction   = 1
	// DotGithubFileTypeWorkflow represents the workflow file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeWorkflow = 2
)

// Glitch represents a linting error.
type Glitch struct {
	Type     int
	Name     string
	Path     string
	RuleName string
	ErrText  string
	IsError  bool
}

// ListToMarkdown takes a list of Glitch instances and generates a Markdown table from it.
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
