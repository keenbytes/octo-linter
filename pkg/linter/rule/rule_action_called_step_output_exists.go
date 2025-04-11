package rule

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.pl/mikogs/octo-linter/pkg/action"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

type RuleActionCalledStepOutputExists struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleActionCalledStepOutputExists) Validate() error {
	return nil
}

func (r RuleActionCalledStepOutputExists) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeAction {
		return
	}
	a := f.(*action.Action)

	if r.Value {
		re := regexp.MustCompile(`\${{[ ]*steps\.([a-zA-Z0-9\-_]+)\.outputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
		reAppendToGithubOutput := regexp.MustCompile(`echo[ ]+["']([a-zA-Z0-9\-_]+)=.*["'][ ]+.*>>[ ]+\$GITHUB_OUTPUT`)
		reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
		reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

		found := re.FindAllSubmatch(a.Raw, -1)
		for _, f := range found {
			stepName := string(f[1])
			outputName := string(f[2])

			if a.Runs == nil {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a step output '%s' but 'runs' does not exist", a.DirName, stepName), chWarnings, chErrors)
				compliant = false
				continue
			}
			step := a.Runs.GetStep(string(f[1]))
			if step == nil {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a step '%s' output '%s' but step does not exist", a.DirName, stepName, outputName), chWarnings, chErrors)
				compliant = false
				continue
			}

			foundOutput := false

			// search in 'run' when there is no 'uses'
			if step.Uses == "" && step.Run != "" {
				foundEchoLines := reAppendToGithubOutput.FindAllSubmatch([]byte(step.Run), -1)
				for _, f := range foundEchoLines {
					if outputName == string(f[1]) {
						foundOutput = true
					}
				}
				if !foundOutput {
					printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a step '%s' output '%s' that does not exist", a.DirName, stepName, outputName), chWarnings, chErrors)
					compliant = false
					continue
				}
			}

			if foundOutput {
				continue
			}

			var action *action.Action
			// local action
			if reLocal.MatchString(step.Uses) {
				actionName := strings.Replace(step.Uses, "./.github/actions/", "", -1)
				action = d.GetAction(actionName)
			}
			// external action
			if reExternal.MatchString(step.Uses) {
				action = d.GetExternalAction(step.Uses)
			}
			if action == nil {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls a step '%s' output '%s' on action that does not exist", a.DirName, stepName, outputName), chWarnings, chErrors)
				compliant = false
				continue
			}

			for duaOutputName := range action.Outputs {
				if duaOutputName == outputName {
					foundOutput = true
				}
			}
			if !foundOutput {
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("action '%s' calls step '%s' output '%s' on action and that output does not exist", a.DirName, stepName, outputName), chWarnings, chErrors)
				compliant = false
				continue
			}
		}
	}

	return
}

func (r RuleActionCalledStepOutputExists) GetConfigName() string {
	return r.ConfigName
}
