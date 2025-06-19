package linter

import (
	"github.com/keenbytes/octo-linter/internal/linter/rule"
	"github.com/keenbytes/octo-linter/internal/linter/rule/filenames"
	"github.com/keenbytes/octo-linter/internal/linter/rule/refvars"
	"github.com/keenbytes/octo-linter/internal/linter/rule/usedactions"
	"github.com/keenbytes/octo-linter/internal/linter/rule/dependencies"
	"github.com/keenbytes/octo-linter/internal/linter/rule/workflowrunners"
)

func (cfg *Config) addRuleFromConfig(fullRuleName string, ruleConfig interface{}) error {
	var ruleInstance rule.Rule

	switch fullRuleName {
	case "dependencies__action_referenced_input_must_exists":
		ruleInstance = dependencies.ActionReferencedInputExists{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "dependencies__action_referenced_step_output_must_exist":
		ruleInstance = dependencies.ActionReferencedStepOutputExists{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "dependencies__workflow_needs_field_must_contain_already_existing_jobs":
		ruleInstance = dependencies.WorkflowNeedsWithExistingJobs{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "dependencies__workflow_referenced_input_must_exists":
		ruleInstance = dependencies.WorkflowReferencedInputExists{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "dependencies__workflow_referenced_variable_must_exists_in_attached_file":
		ruleInstance = dependencies.WorkflowReferencedVariableExistsInFile{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "filenames__action_directory_name_format":
		ruleInstance = filenames.ActionDirectoryNameFormat{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "filenames__action_filename_extensions_allowed":
		ruleInstance = filenames.FilenameExtensionsAllowed{}
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
		ruleInstance = filenames.FilenameExtensionsAllowed{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_actions__not_in_double_quotes":
		ruleInstance = refvars.NotInDoubleQuotes{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_actions__not_one_word":
		ruleInstance = refvars.NotOneWord{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_workflows__not_in_double_quotes":
		ruleInstance = refvars.NotInDoubleQuotes{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "referenced_variables_in_workflows__not_one_word":
		ruleInstance = refvars.NotOneWord{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_action_steps__must_exist":
		ruleInstance = usedactions.Exists{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_action_steps__must_have_valid_inputs":
		ruleInstance = usedactions.ValidInputs{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_action_steps__source":
		ruleInstance = usedactions.Source{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_workflow_job_steps__must_exist":
		ruleInstance = usedactions.Exists{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_workflow_job_steps__must_have_valid_inputs":
		ruleInstance = usedactions.ValidInputs{}
		err := ruleInstance.Validate(ruleConfig)
		if err != nil {
			return err
		}
	case "used_actions_in_workflow_job_steps__source":
		ruleInstance = usedactions.Source{}
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
