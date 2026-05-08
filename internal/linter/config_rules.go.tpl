package linter

import (
	"fmt"

	"octo-linter/internal/linter/rule"
	"octo-linter/internal/linter/rule/filenames"
	"octo-linter/internal/linter/rule/naming"
	"octo-linter/internal/linter/rule/required"
	"octo-linter/internal/linter/rule/refvars"
	"octo-linter/internal/linter/rule/usedactions"
	"octo-linter/internal/linter/rule/dependencies"
	"octo-linter/internal/linter/rule/runners"
)

//nolint:gocognit,gocyclo,funlen,maintidx
func (cfg *Config) addRuleFromConfig(fullRuleName string, ruleConfig interface{}) error {
	var ruleInstance rule.Rule

	switch fullRuleName {

  {{- range $configName, $structDetails := .Rules }}
	case "{{ $configName }}":
		ruleInstance = {{ $structDetails.N }}{
			{{- range $fieldName, $fieldValue := $structDetails.F }}
			{{ $fieldName }}: {{ $fieldValue }},
			{{- end }}
		}

		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return fmt.Errorf("rule validation error: %w", err)
		}
  {{- end }}
	}

	if ruleInstance != nil {
		cfg.Rules = append(cfg.Rules, ruleInstance)
		cfg.Values = append(cfg.Values, ruleConfig)
	}

	return nil
}
