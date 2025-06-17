package linter

import (
	"github.com/keenbytes/octo-linter/pkg/linter/rule"
	rulefilenames "github.com/keenbytes/octo-linter/pkg/linter/rule/filenames"
)

func (cfg *Config) addRuleFromConfig(fullRuleName string, ruleConfig interface{}) error {
	var ruleInstance rule.Rule

	switch fullRuleName {
	case "filenames__action_filename_extensions_allowed":
		ruleInstance = rulefilenames.ActionFilenameExtensionsAllowed{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	default:
		// do nothing for now
	}

	if ruleInstance != nil {
		cfg.Rules = append(cfg.Rules, ruleInstance)
		cfg.Values = append(cfg.Values, ruleConfig)
	}

	return nil
}
