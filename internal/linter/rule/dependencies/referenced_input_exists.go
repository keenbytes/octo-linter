package dependencies

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// ReferencedInputExists scans the code for all input references and verifies that each has been previously defined.
// During action or workflow execution, if a reference to an undefined input is found, it is replaced with an empty string.
type ReferencedInputExists struct {
}

func (r ReferencedInputExists) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "dependencies__workflow_referenced_input_must_exists"
	case rule.DotGithubFileTypeAction:
		return "dependencies__action_referenced_input_must_exists"
	default:
		return "dependencies__*_referenced_input_must_exists"
	}
}

func (r ReferencedInputExists) FileType() int {
	return rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow
}

func (r ReferencedInputExists) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r ReferencedInputExists) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeAction && f.GetType() != rule.DotGithubFileTypeWorkflow {
		return
	}

	if !conf.(bool) {
		return
	}

	if f.GetType() == rule.DotGithubFileTypeAction {
		a := f.(*action.Action)

		re := regexp.MustCompile(`\${{[ ]*inputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			if a.Inputs == nil || a.Inputs[string(f[1])] == nil {
				chErrors <- fmt.Sprintf("action '%s' calls an input '%s' that does not exist", a.DirName, string(f[1]))
				compliant = false
			}
		}
	}

	if f.GetType() == rule.DotGithubFileTypeWorkflow {
		w := f.(*workflow.Workflow)
		re := regexp.MustCompile(`\${{[ ]*inputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
		found := re.FindAllSubmatch(w.Raw, -1)
		for _, f := range found {
			notInInputs := true
			if w.On != nil {
				if w.On.WorkflowCall != nil && w.On.WorkflowCall.Inputs != nil && w.On.WorkflowCall.Inputs[string(f[1])] != nil {
					notInInputs = false
				}
				if w.On.WorkflowDispatch != nil && w.On.WorkflowDispatch.Inputs != nil && w.On.WorkflowDispatch.Inputs[string(f[1])] != nil {
					notInInputs = false
				}
			}
			if notInInputs {
				chErrors <- fmt.Sprintf("workflow '%s' calls an input '%s' that does not exist", w.FileName, string(f[1]))
				compliant = false
			}
		}
	}

	return
}
