package naming

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/casematch"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// Action checks if specified action field adheres to the selected naming convention.
type Action struct {
	Field string
}

func (r Action) ConfigName(int) string {
	switch r.Field {
	case "input_name":
		return "naming_conventions__action_input_name_format"
	case "output_name":
		return "naming_conventions__action_output_name_format"
	case "referenced_variable":
		return "naming_conventions__action_referenced_variable_format"
	case "step_env":
		return "naming_conventions__action_step_env_format"
	default:
		return "naming_conventions__action_*"
	}
}

func (r Action) FileType() int {
	return rule.DotGithubFileTypeAction
}

func (r Action) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != "dash-case" && val != "camelCase" && val != "PascalCase" && val != "ALL_CAPS" {
		return fmt.Errorf("value can be one of: dash-case, camelCase, PascalCase, ALL_CAPS")
	}

	return nil
}

func (r Action) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction {
		return true, nil
	}
	a := f.(*action.Action)

	compliant := true

	switch r.Field {
	case "input_name":
		for inputName := range a.Inputs {
			m := casematch.Match(inputName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path:     a.Path,
					Name:     a.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  fmt.Sprintf("input '%s' must be %s", inputName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	case "output_name":
		for outputName := range a.Outputs {
			m := casematch.Match(outputName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path:     a.Path,
					Name:     a.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  fmt.Sprintf("output '%s' must be %s", outputName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	case "referenced_variable":
		varTypes := []string{"env", "var", "secret"}
		for _, v := range varTypes {
			re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
			found := re.FindAllSubmatch(a.Raw, -1)
			for _, f := range found {
				m := casematch.Match(string(f[1]), conf.(string))
				if !m {
					chErrors <- glitch.Glitch{
						Path:     a.Path,
						Name:     a.DirName,
						Type:     rule.DotGithubFileTypeAction,
						ErrText:  fmt.Sprintf("references a variable '%s' that must be %s", string(f[1]), conf.(string)),
						RuleName: r.ConfigName(0),
					}
					compliant = false
				}
			}
		}
	case "step_env":
		if len(a.Runs.Steps) == 0 {
			return true, nil
		}

		for i, step := range a.Runs.Steps {
			if len(step.Env) == 0 {
				continue
			}
			for envName := range step.Env {
				m := casematch.Match(envName, conf.(string))
				if !m {
					chErrors <- glitch.Glitch{
						Path:     a.Path,
						Name:     a.DirName,
						Type:     rule.DotGithubFileTypeAction,
						ErrText:  fmt.Sprintf("step %d env '%s' must be %s", i, envName, conf.(string)),
						RuleName: r.ConfigName(0),
					}
					compliant = false
				}
			}
		}
	default:
		// do nothing
	}

	return compliant, nil
}
