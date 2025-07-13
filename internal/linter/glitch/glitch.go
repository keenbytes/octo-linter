// Package glitch contains code related to representing a lint error.
package glitch

import (
	"fmt"
)

const (
	// DotGithubFileTypeAction represents the action file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeAction = 1
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
func ListToMarkdown(glitches []*Glitch, limit int) string {
	markdown := `|Item|Error|
|---|---|
`

	for i, glitch := range glitches {
		if limit > 0 && i == limit {
			break
		}

		name := "a/" + glitch.Name
		if glitch.Type == DotGithubFileTypeWorkflow {
			name = "w/" + glitch.Name
		}

		level := `🟠`
		if glitch.IsError {
			level = `🔴`
		}

		markdown += fmt.Sprintf("|%s|%s %s *(%s)*|\n", name, level, glitch.ErrText, glitch.RuleName)
	}

	if len(glitches) > limit && limit > 0 {
		markdown += fmt.Sprintf("\n...and many more (%d in total).", len(glitches))
	}

	return markdown
}
