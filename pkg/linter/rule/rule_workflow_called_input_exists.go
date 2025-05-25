package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

// RuleWorkflowCalledInputExists scans the code for all input references and verifies that each has been previously defined.
// During execution, if a reference to an undefined input is found, it is replaced with an empty string.
type RuleWorkflowCalledInputExists struct {
	Value      bool
	ConfigName string
	IsError    bool
}

func (r RuleWorkflowCalledInputExists) Validate() error {
	return nil
}

func (r RuleWorkflowCalledInputExists) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value {
		return
	}

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
			printErrOrWarn(r.ConfigName, r.IsError, fmt.Sprintf("workflow '%s' calls an input '%s' that does not exist", w.FileName, string(f[1])), chWarnings, chErrors)
			compliant = false
		}
	}

	return
}

func (r RuleWorkflowCalledInputExists) GetConfigName() string {
	return r.ConfigName
}
