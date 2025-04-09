package linter

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
)

const (
	HasNoErrorsOrWarnings = iota
	HasErrors
	HasOnlyWarnings
)

type Linter struct {
	Config   *Config
	LogLevel int
}

func (l *Linter) Lint(d *dotgithub.DotGithub) (uint8, error) {
	if l.Config == nil {
		panic("Config cannot be nil")
	}
	if d == nil {
		panic("DotGithub cannot be empty")
	}

	finalStatus := HasNoErrorsOrWarnings

	for _, action := range d.Actions {
		for _, rule := range l.Config.ActionRules {
			_, isError := l.Config.Errors[rule.GetConfigName()]
			compliant, err := rule.Lint(action, d)
			if err != nil {
				return HasErrors, fmt.Errorf("error with running rule %s: %w", rule.GetConfigName(), err)
			}
			if !compliant {
				if isError {
					finalStatus = HasErrors
				} else {
					if finalStatus == HasNoErrorsOrWarnings {
						finalStatus = HasOnlyWarnings
					}
				}
			}
		}
	}

	for _, workflow := range d.Workflows {
		for _, rule := range l.Config.WorkflowRules {
			_, isError := l.Config.Errors[rule.GetConfigName()]
			compliant, err := rule.Lint(workflow, d)
			if err != nil {
				return HasErrors, fmt.Errorf("error with running rule %s: %w", rule.GetConfigName(), err)
			}
			if !compliant {
				if isError {
					finalStatus = HasErrors
				} else {
					if finalStatus == HasNoErrorsOrWarnings {
						finalStatus = HasOnlyWarnings
					}
				}
			}
		}
	}

	for _, rule := range l.Config.DotGithubRules {
		_, isError := l.Config.Errors[rule.GetConfigName()]
		compliant, err := rule.Lint(d)
		if err != nil {
			return HasErrors, fmt.Errorf("error with running rule %s: %w", rule.GetConfigName(), err)
		}
		if !compliant {
			if isError {
				finalStatus = HasErrors
			} else {
				if finalStatus == HasNoErrorsOrWarnings {
					finalStatus = HasOnlyWarnings
				}
			}
		}
	}

	return uint8(finalStatus), nil
}
