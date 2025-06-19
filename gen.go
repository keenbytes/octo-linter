package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func main() {
	genPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "../../" {
		genPath = filepath.Join(genPath, os.Args[1])
	}

	tplObj := struct {
		Rules map[string]string
	}{
		Rules: map[string]string{
			"filenames__action_filename_extensions_allowed":                           "filenames.ActionFilenameExtensionsAllowed",
			"filenames__action_directory_name_format":                                 "filenames.ActionDirectoryNameFormat",
			"filenames__workflow_filename_extensions_allowed":                         "filenames.WorkflowFilenameExtensionsAllowed",
			"filenames__workflow_filename_base_format":                                "filenames.WorkflowFilenameBaseFormat",
			"workflow_runners__not_latest":                                            "workflowrunners.NotLatest",
			"referenced_variables_in_actions__not_one_word":                           "refvars.NotOneWord_InAction",
			"referenced_variables_in_actions__not_in_double_quotes":                   "refvars.NotInDoubleQuotes_InAction",
			"referenced_variables_in_workflows__not_one_word":                         "refvars.NotOneWord_InWorkflow",
			"referenced_variables_in_workflows__not_in_double_quotes":                 "refvars.NotInDoubleQuotes_InWorkflow",
			"dependencies__workflow_needs_field_must_contain_already_existing_jobs":   "dependencies.WorkflowNeedsWithExistingJobs",
			"dependencies__action_referenced_input_must_exists":                       "dependencies.ActionReferencedInputExists",
			"dependencies__action_referenced_step_output_must_exist":                  "dependencies.ActionReferencedStepOutputExists",
			"dependencies__workflow_referenced_variable_must_exists_in_attached_file": "dependencies.WorkflowReferencedVariableExistsInFile",
			"dependencies__workflow_referenced_input_must_exists":                     "dependencies.WorkflowReferencedInputExists",
			"used_actions_in_action_steps__source":                                    "usedactions.Source",
			"used_actions_in_action_steps__must_exist":                                "usedactions.Exists",
			"used_actions_in_action_steps__must_have_valid_inputs":                    "usedactions.ValidInputs",
			"used_actions_in_workflow_job_steps__source":                              "usedactions.Source",
			"used_actions_in_workflow_job_steps__must_exist":                          "usedactions.Exists",
			"used_actions_in_workflow_job_steps__must_have_valid_inputs":              "usedactions.ValidInputs",
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
