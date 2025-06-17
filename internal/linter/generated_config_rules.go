package linter

import (
	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/internal/linter/rule/filenames"
)

func (cfg *Config) addRuleFromConfig(fullRuleName string, ruleConfig interface{}) error {
	var ruleInstance rule.Rule

	switch fullRuleName {
	case "filenames__action_directory_name_format":
		ruleInstance = filenames.ActionDirectoryNameFormat{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "filenames__action_filename_extensions_allowed":
		ruleInstance = filenames.ActionFilenameExtensionsAllowed{}
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
