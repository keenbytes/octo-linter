package linter

import (
	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/internal/linter/rule/filenames"
	"github.com/keenbytes/octo-linter/internal/linter/rule/workflowrunners"
	"github.com/keenbytes/octo-linter/internal/linter/rule/refvars"
)

func (cfg *Config) addRuleFromConfig(fullRuleName string, ruleConfig interface{}) error {
	var ruleInstance rule.Rule

	switch fullRuleName {

  {{- range $configName, $structName := .Rules }}
	case "{{ $configName }}":
		ruleInstance = {{ $structName }}{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
  {{- end }}

	default:
		// do nothing for now
	}

	if ruleInstance != nil {
		cfg.Rules = append(cfg.Rules, ruleInstance)
		cfg.Values = append(cfg.Values, ruleConfig)
	}

	return nil
}
