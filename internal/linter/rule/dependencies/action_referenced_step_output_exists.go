package dependencies

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

// ActionReferencedStepOutputExists checks whether references to step outputs correspond to outputs defined in preceding steps.
// During execution, referencing a non-existent step output results in an empty string.
type ActionReferencedStepOutputExists struct{}

// ConfigName returns the name of the rule as defined in the configuration file.
func (r ActionReferencedStepOutputExists) ConfigName(int) string {
	return "dependencies__action_referenced_step_output_must_exist"
}

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r ActionReferencedStepOutputExists) FileType() int {
	return rule.DotGithubFileTypeAction
}

// Validate checks whether the given value is valid for this rule's configuration.
func (r ActionReferencedStepOutputExists) Validate(conf interface{}) error {
	_, ok := conf.(bool)
	if !ok {
		return errors.New("value should be bool")
	}

	return nil
}

// Lint runs a rule with the specified configuration on a dotgithub.File (action or workflow),
// reports any errors via the given channel, and returns whether the file is compliant.
func (r ActionReferencedStepOutputExists) Lint(
	conf interface{},
	f dotgithub.File,
	d *dotgithub.DotGithub,
	chErrors chan<- glitch.Glitch,
) (bool, error) {
	err := r.Validate(conf)
	if err != nil {
		return false, err
	}

	if f.GetType() != rule.DotGithubFileTypeAction || !conf.(bool) {
		return true, nil
	}

	a := f.(*action.Action)

	re := regexp.MustCompile(`\${{[ ]*steps\.([a-zA-Z0-9\-_]+)\.outputs\.([a-zA-Z0-9\-_]+)[ ]*}}`)
	reAppendToGithubOutput := regexp.MustCompile(
		`echo[ ]+["']([a-zA-Z0-9\-_]+)=.*["'][ ]+.*>>[ ]+\$GITHUB_OUTPUT`,
	)
	reLocal := regexp.MustCompile(`^\.\/\.github\/actions\/([a-z0-9\-]+|[a-z0-9\-]+\/[a-z0-9\-]+)$`)
	reExternal := regexp.MustCompile(
		`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`,
	)

	compliant := true

	found := re.FindAllSubmatch(a.Raw, -1)
	for _, f := range found {
		stepName := string(f[1])
		outputName := string(f[2])

		if a.Runs == nil {
			chErrors <- glitch.Glitch{
				Path:     a.Path,
				Name:     a.DirName,
				Type:     rule.DotGithubFileTypeAction,
				ErrText:  fmt.Sprintf("calls a step output '%s' but 'runs' does not exist", stepName),
				RuleName: r.ConfigName(0),
			}

			compliant = false

			continue
		}

		step := a.Runs.GetStep(string(f[1]))
		if step == nil {
			chErrors <- glitch.Glitch{
				Path:     a.Path,
				Name:     a.DirName,
				Type:     rule.DotGithubFileTypeAction,
				ErrText:  fmt.Sprintf("calls a step '%s' output '%s' but step does not exist", stepName, outputName),
				RuleName: r.ConfigName(0),
			}

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
				chErrors <- glitch.Glitch{
					Path:     a.Path,
					Name:     a.DirName,
					Type:     rule.DotGithubFileTypeAction,
					ErrText:  fmt.Sprintf("calls a step '%s' output '%s' that does not exist", stepName, outputName),
					RuleName: r.ConfigName(0),
				}

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
			actionName := strings.ReplaceAll(step.Uses, "./.github/actions/", "")
			action = d.GetAction(actionName)
		}
		// external action
		if reExternal.MatchString(step.Uses) {
			action = d.GetExternalAction(step.Uses)
		}

		if action == nil {
			chErrors <- glitch.Glitch{
				Path:     a.Path,
				Name:     a.DirName,
				Type:     rule.DotGithubFileTypeAction,
				ErrText:  fmt.Sprintf("calls a step '%s' output '%s' on action that does not exist", stepName, outputName),
				RuleName: r.ConfigName(0),
			}

			compliant = false

			continue
		}

		for duaOutputName := range action.Outputs {
			if duaOutputName == outputName {
				foundOutput = true
			}
		}

		if !foundOutput {
			chErrors <- glitch.Glitch{
				Path:     a.Path,
				Name:     a.DirName,
				Type:     rule.DotGithubFileTypeAction,
				ErrText:  fmt.Sprintf("calls step '%s' output '%s' on action and that output does not exist", stepName, outputName),
				RuleName: r.ConfigName(0),
			}

			compliant = false

			continue
		}
	}

	return compliant, nil
}
