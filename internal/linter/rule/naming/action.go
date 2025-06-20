package naming

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/casematch"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
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

func (r Action) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	switch r.Field {
	case "input_name":
		for inputName := range a.Inputs {
			m := casematch.Match(inputName, conf.(string))
			if !m {
				chErrors <- fmt.Sprintf("action '%s' input '%s' must be %s", a.DirName, inputName, conf.(string))
				compliant = false
			}
		}
	case "output_name":
		for outputName := range a.Outputs {
			m := casematch.Match(outputName, conf.(string))
			if !m {
				chErrors <- fmt.Sprintf("action '%s' output '%s' must be %s", a.DirName, outputName, conf.(string))
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
					chErrors <- fmt.Sprintf("action '%s' references a variable '%s' that must be %s", a.DirName, string(f[1]), conf.(string))
					compliant = false
				}
			}
		}
	case "step_env":
		if a.Runs == nil || a.Runs.Steps == nil || len(a.Runs.Steps) == 0 {
			return
		}

		for i, step := range a.Runs.Steps {
			if step.Env == nil || len(step.Env) == 0 {
				continue
			}
			for envName := range step.Env {
				m := casematch.Match(envName, conf.(string))
				if !m {
					chErrors <- fmt.Sprintf("action '%s' step %d env '%s' must be %s", a.DirName, i, envName, conf.(string))
					compliant = false
				}
			}
		}
	default:
		// do nothing
	}

	return
}
