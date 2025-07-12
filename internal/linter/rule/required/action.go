package required

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// Action checks if required fields within actions are defined.
type Action struct {
	Field int
}

const (
	_ = iota
	// ActionFieldAction specifies that the rule targets top-level fields in a GitHub Actions action.
	ActionFieldAction
	// ActionFieldInput specifies that the rule targets the 'inputs' section.
	ActionFieldInput
	// ActionFieldOutput specifies that the rule targets the 'outputs' section.
	ActionFieldOutput
)

// ConfigName returns the name of the rule as defined in the configuration file.
func (r Action) ConfigName(int) string {
	switch r.Field {
	case ActionFieldAction:
		return "required_fields__action_requires"
	case ActionFieldInput:
		return "required_fields__action_input_requires"
	case ActionFieldOutput:
		return "required_fields__action_output_requires"
	default:
		return "required_fields__action_*_requires"
	}
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r Action) FileType() int {
	return rule.DotGithubFileTypeAction
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r Action) Validate(conf interface{}) error {
	vals, ok := conf.([]interface{})
	if !ok {
		return errors.New("value should be []string")
	}

	for _, v := range vals {
		field, ok := v.(string)
		if !ok {
			return errors.New("value should be []string")
		}

		switch r.Field {
		case ActionFieldAction:
			if field != "name" && field != "description" {
				return errors.New("value can contain only 'name' and/or 'description'")
			}
		case ActionFieldInput, ActionFieldOutput:
			if field != "description" {
				return errors.New("value can contain only 'description'")
			}
		}
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r Action) Lint(
	conf interface{},
	f dotgithub.File,
	_ *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction {
		return true, nil
	}

	a := f.(*action.Action)

	confInterfaces := conf.([]interface{})

	compliant := true

	switch r.Field {
	case ActionFieldAction:
		for _, field := range confInterfaces {
			if (field.(string) == "name" && a.Name == "") ||
				(field.(string) == "description" && a.Description == "") {
				chErrors <- glitch.Glitch{
					Path:     a.Path,
					Name:     a.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  "does not have a required " + field.(string),
					RuleName: r.ConfigName(0),
				}

				compliant = false
			}
		}
	case ActionFieldInput:
		for inputName, input := range a.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- glitch.Glitch{
						Path:     a.Path,
						Name:     a.DirName,
						Type:     rule.DotGithubFileTypeAction,
						ErrText:  fmt.Sprintf("input '%s' does not have a required %s", inputName, field.(string)),
						RuleName: r.ConfigName(0),
					}

					compliant = false
				}
			}
		}
	case ActionFieldOutput:
		for outputName, output := range a.Outputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && output.Description == "" {
					chErrors <- glitch.Glitch{
						Path:     a.Path,
						Name:     a.DirName,
						Type:     rule.DotGithubFileTypeAction,
						ErrText:  fmt.Sprintf("output '%s' does not have a required %s", outputName, field.(string)),
						RuleName: r.ConfigName(0),
					}

					compliant = false
				}
			}
		}
	}

	return compliant, nil
}
