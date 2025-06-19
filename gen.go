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
			"filenames__action_filename_extensions_allowed":   "filenames.ActionFilenameExtensionsAllowed",
			"filenames__action_directory_name_format":         "filenames.ActionDirectoryNameFormat",
			"filenames__workflow_filename_extensions_allowed": "filenames.WorkflowFilenameExtensionsAllowed",
			"filenames__workflow_filename_base_format":        "filenames.WorkflowFilenameBaseFormat",
			"workflow_runners__not_latest":                    "workflowrunners.NotLatest",
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
