package dependencies

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

// WorkflowReferencedInputExists scans the workflow code for all input references and verifies that each has been previously defined.
// During workflow execution, if a reference to an undefined input is found, it is replaced with an empty string.
type WorkflowReferencedInputExists struct {
}

func (r WorkflowReferencedInputExists) ConfigName() string {
	return "dependencies__workflow_referenced_input_must_exists"
}

func (r WorkflowReferencedInputExists) FileType() int {
	return rule.DotGithubFileTypeWorkflow
}

func (r WorkflowReferencedInputExists) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

func (r WorkflowReferencedInputExists) Lint(conf interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != rule.DotGithubFileTypeWorkflow || !conf.(bool) {
		return
	}
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


	return
}
