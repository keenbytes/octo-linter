package linter

import (
	"fmt"

	"gopkg.in/yaml.v2"
	"github.com/keenbytes/octo-linter/pkg/linter/rule"
)

func (cfg *Config) readBytesAndValidate(b []byte) error {
	cfg.Rules = make([]rule.Rule, 0)

	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling: %w", err)
	}

	// Format list of strings into map for config version > 1
	if cfg.Version != "1" {
		cfg.WarningOnlyMap = make(map[string]struct{})
		for _, w := range cfg.WarningOnly {
			cfg.WarningOnlyMap[w] = struct{}{}
		}
	}

	for ruleName, ruleValue := range cfg.RulesConfig {
		isError := cfg.IsError(ruleName)

		switch ruleName {
		case "action_file_extensions":
			cfg.Rules = append(cfg.Rules, rule.RuleActionFileExtensions{
				Value:      iArrToStrArr(ruleValue),
				ConfigName: "action_file_extensions",
				IsError:    isError,
			})
		case "action_directory_name":
			cfg.Rules = append(cfg.Rules, rule.RuleActionDirectoryName{
				Value:      ruleValue.(string),
				ConfigName: "action_directory_name",
				IsError:    isError,
			})
		case "workflow_file_extensions":
			cfg.Rules = append(cfg.Rules, rule.RuleWorkflowFileExtensions{
				Value:      iArrToStrArr(ruleValue),
				ConfigName: "workflow_file_extensions",
				IsError:    isError,
			})
		case "action_called_variable":
			cfg.Rules = append(cfg.Rules, rule.RuleActionCalledVariable{
				Value:      ruleValue.(string),
				ConfigName: "action_called_variable",
				IsError:    isError,
			})
		case "action_called_variable_not_one_word":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleActionCalledVariableNotOneWord{
					Value:      true,
					ConfigName: "action_called_variable_not_one_word",
					IsError:    isError,
				})
			}
		case "action_called_variable_not_in_double_quote":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleActionCalledVariableNotInDoubleQuote{
					Value:      true,
					ConfigName: "action_called_variable_not_in_double_quote",
					IsError:    isError,
				})
			}
		case "action_called_input_exists":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleActionCalledInputExists{
					Value:      true,
					ConfigName: "action_called_input_exists",
					IsError:    isError,
				})
			}
		case "action_called_step_output_exists":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleActionCalledStepOutputExists{
					Value:      true,
					ConfigName: "action_called_step_output_exists",
					IsError:    isError,
				})
			}
		case "action_step_action":
			cfg.Rules = append(cfg.Rules, rule.RuleStepAction{
				Value:      ruleValue.(string),
				ConfigName: "action_step_action",
				IsError:    isError,
			})
		case "action_step_action_input_valid":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleStepActionInputValid{
					Value:      true,
					ConfigName: "action_step_action_input_valid",
					IsError:    isError,
				})
			}
		case "action_step_env":
			cfg.Rules = append(cfg.Rules, rule.RuleStepEnv{
				Value:      ruleValue.(string),
				ConfigName: "action_step_env",
				IsError:    isError,
			})
		case "workflow_env":
			cfg.Rules = append(cfg.Rules, rule.RuleWorkflowEnv{
				Value:      ruleValue.(string),
				ConfigName: "workflow_env",
				IsError:    isError,
			})
		case "workflow_called_variable":
			cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCalledVariable{
				Value:      ruleValue.(string),
				ConfigName: "workflow_called_variable",
				IsError:    isError,
			})
		case "workflow_called_variable_not_one_word":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCalledVariableNotOneWord{
					Value:      true,
					ConfigName: "workflow_called_variable_not_one_word",
					IsError:    isError,
				})
			}
		case "workflow_called_variable_not_in_double_quote":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCalledVariableNotInDoubleQuote{
					Value:      true,
					ConfigName: "workflow_called_variable_not_in_double_quote",
					IsError:    isError,
				})
			}
		case "workflow_called_variable_exists_in_file":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCalledVariableExistsInFile{
					Value:      true,
					ConfigName: "workflow_called_variable_exists_in_file",
					IsError:    isError,
				})
			}
		case "workflow_called_input_exists":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCalledInputExists{
					Value:      true,
					ConfigName: "workflow_called_input_exists",
					IsError:    isError,
				})
			}
		case "workflow_single_job_main":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowSingleJobMain{
					Value:      true,
					ConfigName: "workflow_single_job_main",
					IsError:    isError,
				})
			}
		case "workflow_job_needs_exist":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowJobNeedsExist{
					Value:      true,
					ConfigName: "workflow_job_needs_exist",
					IsError:    isError,
				})
			}
		case "workflow_required_uses_or_runs_on":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowRequiredUsesOrRunsOn{
					Value:      true,
					ConfigName: "workflow_required_uses_or_runs_on",
					IsError:    isError,
				})
			}
		case "workflow_runs_on_not_latest":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleWorkflowRunsOnNotLatest{
					Value:      true,
					ConfigName: "workflow_runs_on_not_latest",
					IsError:    isError,
				})
			}
		case "workflow_job_env":
			cfg.Rules = append(cfg.Rules, rule.RuleWorkflowJobEnv{
				Value:      ruleValue.(string),
				ConfigName: "workflow_job_env",
				IsError:    isError,
			})
		case "workflow_job_step_action":
			cfg.Rules = append(cfg.Rules, rule.RuleStepAction{
				Value:      ruleValue.(string),
				ConfigName: "workflow_job_step_action",
				IsError:    isError,
			})
		case "workflow_job_step_action_input_valid":
			if ruleValue.(bool) {
				cfg.Rules = append(cfg.Rules, rule.RuleStepActionInputValid{
					Value:      true,
					ConfigName: "workflow_job_step_action_input_valid",
					IsError:    isError,
				})
			}
		case "workflow_job_step_env":
			cfg.Rules = append(cfg.Rules, rule.RuleStepEnv{
				Value:      ruleValue.(string),
				ConfigName: "workflow_job_step_env",
				IsError:    isError,
			})
		case "action_required__name", "action_required__description":
		case "action_input_required__description", "action_input_value__name":
		case "action_output_required__description", "action_output_value__name":
		case "action_step_action_exists__local", "action_step_action_exists__external":
		case "workflow_required__name":
		case "workflow_dispatch_input_required__description", "workflow_dispatch_input_value__name":
		case "workflow_call_input_required__description", "workflow_call_input_value__name":
		case "workflow_job_value__name":
		case "workflow_job_step_action_exists__local", "workflow_job_step_action_exists__external":
		default:
			return fmt.Errorf("invalid rule %s", ruleName)
		}
	}

	cfg.addActionRequired()
	cfg.addActionInputRequired()
	cfg.addActionOutputRequired()
	cfg.addActionInputValue()
	cfg.addActionOutputValue()
	cfg.addActionStepActionExists()
	cfg.addWorkflowRequired()
	cfg.addWorkflowDispatchInputValue()
	cfg.addWorkflowCallInputValue()
	cfg.addWorkflowDispatchInputRequired()
	cfg.addWorkflowCallInputRequired()
	cfg.addWorkflowJobValue()
	cfg.addWorkflowJobStepActionExists()

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("errors have been found: %w", err)
	}

	return nil
}

func (cfg *Config) addActionRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("action_required", []string{"name", "description"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleActionRequired{
			Value:      ruleValue,
			ConfigName: "action_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addActionInputRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("action_input_required", []string{"description"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleActionInputRequired{
			Value:      ruleValue,
			ConfigName: "action_input_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addActionOutputRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("action_output_required", []string{"description"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleActionOutputRequired{
			Value:      ruleValue,
			ConfigName: "action_output_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addActionInputValue() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesWithMapValueIntoOne("action_input_value", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleActionInputValue{
			Value:      ruleValue,
			ConfigName: "action_input_value",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addActionOutputValue() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesWithMapValueIntoOne("action_output_value", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleActionOutputValue{
			Value:      ruleValue,
			ConfigName: "action_output_value",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addActionStepActionExists() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("action_step_action_exists", []string{"local", "external"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleStepActionExists{
			Value:      ruleValue,
			ConfigName: "action_step_action_exists",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("workflow_required", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowRequired{
			Value:      ruleValue,
			ConfigName: "workflow_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowCallInputRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("workflow_call_input_required", []string{"description"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCallInputRequired{
			Value:      ruleValue,
			ConfigName: "workflow_call_input_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowDispatchInputRequired() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("workflow_dispatch_input_required", []string{"description"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowDispatchInputRequired{
			Value:      ruleValue,
			ConfigName: "workflow_dispatch_input_required",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowCallInputValue() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesWithMapValueIntoOne("workflow_call_input_value", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowCallInputValue{
			Value:      ruleValue,
			ConfigName: "workflow_call_input_value",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowDispatchInputValue() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesWithMapValueIntoOne("workflow_dispatch_input_value", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowDispatchInputValue{
			Value:      ruleValue,
			ConfigName: "workflow_dispatch_input_value",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowJobValue() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesWithMapValueIntoOne("workflow_job_value", []string{"name"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleWorkflowJobValue{
			Value:      ruleValue,
			ConfigName: "workflow_job_value",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) addWorkflowJobStepActionExists() {
	ruleValue, ruleIsError := cfg.mergeMultipleRulesIntoOne("workflow_job_step_action_exists", []string{"local", "external"})
	if len(ruleValue) > 0 {
		cfg.Rules = append(cfg.Rules, rule.RuleStepActionExists{
			Value:      ruleValue,
			ConfigName: "workflow_job_step_action_exists",
			IsError:    ruleIsError,
		})
	}
}

func (cfg *Config) mergeMultipleRulesIntoOne(rulePrefix string, ruleVariants []string) ([]string, []bool) {
	ruleValue := make([]string, 0, len(ruleVariants))
	ruleIsError := make([]bool, 0, len(ruleVariants))
	for _, variant := range ruleVariants {
		ruleVariant := fmt.Sprintf("%s__%s", rulePrefix, variant)
		val, ex := cfg.RulesConfig[ruleVariant]
		exErr := cfg.IsError(ruleVariant)
		if ex && val.(bool) {
			ruleValue = append(ruleValue, variant)
			if exErr {
				ruleIsError = append(ruleIsError, true)
			} else {
				ruleIsError = append(ruleIsError, false)
			}
		}
	}
	return ruleValue, ruleIsError
}

func (cfg *Config) mergeMultipleRulesWithMapValueIntoOne(rulePrefix string, ruleVariants []string) (map[string]string, map[string]bool) {
	ruleValue := make(map[string]string, len(ruleVariants))
	ruleIsError := make(map[string]bool, len(ruleVariants))
	for _, variant := range ruleVariants {
		ruleVariant := fmt.Sprintf("%s__%s", rulePrefix, variant)
		val, ex := cfg.RulesConfig[ruleVariant]
		exErr := cfg.IsError(ruleVariant)
		if ex && val.(string) != "" {
			ruleValue[variant] = val.(string)
			ruleIsError[variant] = exErr
		}
	}
	return ruleValue, ruleIsError
}

func iArrToStrArr(i interface{}) []string {
	s := make([]string, 0, len(i.([]interface{})))
	for _, iv := range i.([]interface{}) {
		s = append(s, iv.(string))
	}
	return s
}
