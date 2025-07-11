package linter

import (
	"fmt"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/filenames"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/naming"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/required"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/refvars"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/usedactions"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/dependencies"
	"github.com/keenbytes/octo-linter/v2/internal/linter/rule/runners"
)

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
