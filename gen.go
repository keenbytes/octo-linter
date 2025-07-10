package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type S struct {
	N string
	F map[string]string
}

func main() {
	genPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "../../" {
		genPath = filepath.Join(genPath, os.Args[1])
	}

	tplObj := struct {
		Rules map[string]S
	}{
		Rules: map[string]S{
			"filenames__action_filename_extensions_allowed":                           {N: "filenames.FilenameExtensionsAllowed"},
			"filenames__action_directory_name_format":                                 {N: "filenames.ActionDirectoryNameFormat"},
			"filenames__workflow_filename_extensions_allowed":                         {N: "filenames.FilenameExtensionsAllowed"},
			"filenames__workflow_filename_base_format":                                {N: "filenames.WorkflowFilenameBaseFormat"},
			"workflow_runners__not_latest":                                            {N: "runners.NotLatest"},
			"referenced_variables_in_actions__not_one_word":                           {N: "refvars.NotOneWord"},
			"referenced_variables_in_actions__not_in_double_quotes":                   {N: "refvars.NotInDoubleQuotes"},
			"referenced_variables_in_workflows__not_one_word":                         {N: "refvars.NotOneWord"},
			"referenced_variables_in_workflows__not_in_double_quotes":                 {N: "refvars.NotInDoubleQuotes"},
			"dependencies__workflow_needs_field_must_contain_already_existing_jobs":   {N: "dependencies.WorkflowNeedsWithExistingJobs"},
			"dependencies__action_referenced_input_must_exists":                       {N: "dependencies.ReferencedInputExists"},
			"dependencies__action_referenced_step_output_must_exist":                  {N: "dependencies.ActionReferencedStepOutputExists"},
			"dependencies__workflow_referenced_variable_must_exists_in_attached_file": {N: "dependencies.WorkflowReferencedVariableExistsInFile"},
			"dependencies__workflow_referenced_input_must_exists":                     {N: "dependencies.ReferencedInputExists"},
			"used_actions_in_action_steps__source":                                    {N: "usedactions.Source"},
			"used_actions_in_action_steps__must_exist":                                {N: "usedactions.Exists"},
			"used_actions_in_action_steps__must_have_valid_inputs":                    {N: "usedactions.ValidInputs"},
			"used_actions_in_workflow_job_steps__source":                              {N: "usedactions.Source"},
			"used_actions_in_workflow_job_steps__must_exist":                          {N: "usedactions.Exists"},
			"used_actions_in_workflow_job_steps__must_have_valid_inputs":              {N: "usedactions.ValidInputs"},
			"naming_conventions__action_input_name_format":                            {N: "naming.Action", F: map[string]string{"Field": `"input_name"`}},
			"naming_conventions__action_output_name_format":                           {N: "naming.Action", F: map[string]string{"Field": `"output_name"`}},
			"naming_conventions__action_referenced_variable_format":                   {N: "naming.Action", F: map[string]string{"Field": `"referenced_variable"`}},
			"naming_conventions__action_step_env_format":                              {N: "naming.Action", F: map[string]string{"Field": `"step_env"`}},
			"naming_conventions__workflow_env_format":                                 {N: "naming.Workflow", F: map[string]string{"Field": `"env"`}},
			"naming_conventions__workflow_job_env_format":                             {N: "naming.Workflow", F: map[string]string{"Field": `"job_env"`}},
			"naming_conventions__workflow_job_step_env_format":                        {N: "naming.Workflow", F: map[string]string{"Field": `"job_step_env"`}},
			"naming_conventions__workflow_referenced_variable_format":                 {N: "naming.Workflow", F: map[string]string{"Field": `"referenced_variable"`}},
			"naming_conventions__workflow_dispatch_input_name_format":                 {N: "naming.Workflow", F: map[string]string{"Field": `"dispatch_input_name"`}},
			"naming_conventions__workflow_call_input_name_format":                     {N: "naming.Workflow", F: map[string]string{"Field": `"call_input_name"`}},
			"naming_conventions__workflow_job_name_format":                            {N: "naming.Workflow", F: map[string]string{"Field": `"job_name"`}},
			"naming_conventions__workflow_single_job_only_name":                       {N: "naming.WorkflowSingleJobOnlyName"},
			"required_fields__action_requires":                                        {N: "required.Action", F: map[string]string{"Field": `"action"`}},
			"required_fields__action_input_requires":                                  {N: "required.Action", F: map[string]string{"Field": `"input"`}},
			"required_fields__action_output_requires":                                 {N: "required.Action", F: map[string]string{"Field": `"output"`}},
			"required_fields__workflow_requires":                                      {N: "required.Workflow", F: map[string]string{"Field": `"workflow"`}},
			"required_fields__workflow_dispatch_input_requires":                       {N: "required.Workflow", F: map[string]string{"Field": `"dispatch_input"`}},
			"required_fields__workflow_call_input_requires":                           {N: "required.Workflow", F: map[string]string{"Field": `"call_input"`}},
			"required_fields__workflow_requires_uses_or_runs_on_required":             {N: "required.WorkflowUsesOrRunsOn"},
		},
	}

	tpl, err := os.ReadFile(filepath.Join(genPath, "internal", "linter", "config_rules.go.tpl"))
	if err != nil {
		panic(fmt.Sprintf("error opening template file: %s", err.Error()))
	}

	f, err := os.OpenFile(filepath.Join(genPath, "internal", "linter", "generated_config_rules.go"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("error opening file to write to: %s", err.Error()))
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	t := template.Must(template.New("gend_tpl").Parse(string(tpl)))
	err = t.Execute(buf, &tplObj)
	if err != nil {
		panic(fmt.Sprintf("error executing template: %s", err.Error()))
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("error writing generated template: %s", err.Error()))
	}

	log.Printf("Generated %s", filepath.Join(genPath, "internal", "linter", "generated_config_rules.go"))
}
