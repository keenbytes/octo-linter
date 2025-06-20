package required

import (
	"errors"
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// Action checks if required fields within actions are defined
type Action struct {
	Field string
}

func (r Action) ConfigName(int) string {
	switch r.Field {
	case "action":
		return "required_fields__action_requires"
	case "input":
		return "required_fields__action_input_requires"
	case "output":
		return "required_fields__action_output_requires"
	default:
		return "required_fields__action_*_requires"
	}
}

func (r Action) FileType() int {
	return rule.DotGithubFileTypeAction
}

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
		case "action":
			if field != "description" {
				return fmt.Errorf("value can contain only 'description'")
			}
		case "input", "output":
			if field != "name" && field != "description" {
				return fmt.Errorf("value can contain only 'name' and/or 'description'")
			}
		default:
			// nothing
		}
	}

	return nil
}

func (r Action) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	confInterfaces := conf.([]interface{})

	switch r.Field {
	case "action":
		for _, field := range confInterfaces {
			if (field.(string) == "name" && a.Name == "") || (field.(string) == "description" && a.Description == "") {
				chErrors <- fmt.Sprintf("action '%s' does not have a required %s", a.DirName, field.(string))
				compliant = false
			}
		}
	case "input":
		for inputName, input := range a.Inputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && input.Description == "" {
					chErrors <- fmt.Sprintf("action '%s' input '%s' does not have a required %s", a.DirName, inputName, field.(string))
					compliant = false
				}
			}
		}
	case "output":
		for outputName, output := range a.Outputs {
			for _, field := range confInterfaces {
				if field.(string) == "description" && output.Description == "" {
					chErrors <- fmt.Sprintf("action '%s' output '%s' does not have a required %s", a.DirName, outputName, field.(string))
					compliant = false
				}
			}
		}
	default:
		// do nothing
	}

	return
}
