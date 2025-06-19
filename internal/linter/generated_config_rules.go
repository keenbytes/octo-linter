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
	case "filenames__workflow_filename_base_format":
		ruleInstance = filenames.WorkflowFilenameBaseFormat{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "filenames__workflow_filename_extensions_allowed":
		ruleInstance = filenames.WorkflowFilenameExtensionsAllowed{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_actions__not_in_double_quotes":
		ruleInstance = refvars.ActionNotInDoubleQuotes{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_actions__not_one_word":
		ruleInstance = refvars.ActionNotOneWord{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_workflows__not_in_double_quotes":
		ruleInstance = refvars.WorkflowNotInDoubleQuotes{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_workflows__not_one_word":
		ruleInstance = refvars.WorkflowNotOneWord{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "workflow_runners__not_latest":
		ruleInstance = workflowrunners.NotLatest{}
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
