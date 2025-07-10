package naming

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/casematch"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// Workflow checks if specified workflow field adheres to the selected naming convention.
type Workflow struct {
	Field string
}

func (r Workflow) ConfigName(int) string {
	switch r.Field {
	case "env":
		return "naming_conventions__workflow_env_format"
	case "job_env":
		return "naming_conventions__workflow_job_env_format"
	case "job_step_env":
		return "naming_conventions__workflow_job_step_env_format"
	case "referenced_variable":
		return "naming_conventions__workflow_referenced_variable_format"
	case "dispatch_input_name":
		return "naming_conventions__workflow_dispatch_input_name_format"
	case "call_input_name":
		return "naming_conventions__workflow_call_input_name_format"
	case "job_name":
		return "naming_conventions__workflow_job_name_format"
	default:
		return "naming_conventions__workflow_*"
	}
}

func (r Workflow) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r Workflow) Validate(conf interface{}) error {
	val, ok := conf.(string)
	if !ok {
		return errors.New("value should be string")
	}

	if val != "dash-case" && val != "camelCase" && val != "PascalCase" && val != "ALL_CAPS" {
		return fmt.Errorf("value can be one of: dash-case, camelCase, PascalCase, ALL_CAPS")
	}

	return nil
}

func (r Workflow) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- glitch.Glitch) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	switch r.Field {
	case "env":
		if w.Env == nil || len(w.Env) == 0 {
			return
		}

		for envName := range w.Env {
			m := casematch.Match(envName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText:fmt.Sprintf("env '%s' must be %s", envName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
			}
		}
	case "job_env":
		if w.Jobs == nil || len(w.Jobs) == 0 {
			return
		}

		for jobName, job := range w.Jobs {
			if job.Env == nil || len(job.Env) == 0 {
				continue
			}
			for envName := range job.Env {
				m := casematch.Match(envName, conf.(string))
				if !m {
					chErrors <- glitch.Glitch{
						Path: w.Path,
						Name: w.DisplayName,
						Type: rule.DotGithubFileTypeWorkflow,
						ErrText: fmt.Sprintf("job '%s' env '%s' must be %s", jobName, envName, conf.(string)),
						RuleName: r.ConfigName(0),
					}
				}
			}
		}
	case "job_step_env":
		for jobName, job := range w.Jobs {
			for i, step := range job.Steps {
				if step.Env == nil || len(step.Env) == 0 {
					continue
				}
				for envName := range step.Env {
					m := casematch.Match(envName, conf.(string))
					if !m {
						chErrors <- glitch.Glitch{
							Path: w.Path,
							Name: w.DisplayName,
							Type: rule.DotGithubFileTypeWorkflow,
							ErrText: fmt.Sprintf("job '%s' step %d env '%s' must be %s", jobName, i, envName, conf.(string)),
							RuleName: r.ConfigName(0),
						}
						compliant = false
					}
				}
			}
		}
	case "referenced_variable":
		varTypes := []string{"env", "vars", "secrets"}
		for _, v := range varTypes {
			re := regexp.MustCompile(fmt.Sprintf("\\${{[ ]*%s\\.([a-zA-Z0-9\\-_]+)[ ]*}}", v))
			found := re.FindAllSubmatch(w.Raw, -1)
			for _, f := range found {
				m := casematch.Match(string(f[1]), conf.(string))
				if !m {
					chErrors <- glitch.Glitch{
						Path: w.Path,
						Name: w.DisplayName,
						Type: rule.DotGithubFileTypeWorkflow,
						ErrText: fmt.Sprintf("calls a variable '%s' that must be %s", string(f[1]), conf.(string)),
						RuleName: r.ConfigName(0),
					}
					compliant = false
				}
			}
		}
	case "dispatch_input_name":
		if w.On == nil || w.On.WorkflowDispatch == nil || w.On.WorkflowDispatch.Inputs == nil || len(w.On.WorkflowDispatch.Inputs) == 0 {
			return
		}

		for inputName := range w.On.WorkflowDispatch.Inputs {
			m := casematch.Match(inputName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText: fmt.Sprintf("call input '%s' name must be %s", inputName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	case "call_input_name":
		if w.On == nil || w.On.WorkflowCall == nil || w.On.WorkflowCall.Inputs == nil || len(w.On.WorkflowCall.Inputs) == 0 {
			return
		}

		for inputName := range w.On.WorkflowCall.Inputs {
			m := casematch.Match(inputName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText: fmt.Sprintf("dispatch input '%s' name must be %s", inputName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	case "job_name":
		if w.Jobs == nil || len(w.Jobs) == 0 {
			return
		}

		for jobName := range w.Jobs {
			m := casematch.Match(jobName, conf.(string))
			if !m {
				chErrors <- glitch.Glitch{
					Path: w.Path,
					Name: w.DisplayName,
					Type: rule.DotGithubFileTypeWorkflow,
					ErrText: fmt.Sprintf("job '%s' name must be %s", jobName, conf.(string)),
					RuleName: r.ConfigName(0),
				}
				compliant = false
			}
		}
	default:
		// do nothing
	}
	return
}
