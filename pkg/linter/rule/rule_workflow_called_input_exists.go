package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowCalledInputExists struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowCalledInputExists) Validate() error {
	return nil
}

func (r RuleWorkflowCalledInputExists) Lint(w *workflow.Workflow, d *dotgithub.DotGithub) (compliant bool, err error) {
	compliant = true

	if r.Value {
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
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' calls an input '%s' that does not exist", w.FileName, string(f[1])))
				compliant = false
			}
		}
	}

	return
}

func (r RuleWorkflowCalledInputExists) GetConfigName() string {
	return r.ConfigName
}
