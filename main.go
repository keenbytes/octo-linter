package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/linter"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
	"gopkg.pl/phings/broccli/v2"
)

func main() {
	cli := broccli.NewCLI("octo-linter", "Validates GitHub Actions workflow and action YAML files", "mg@computerclub.pl")

	cmd := cli.AddCmd("lint", "Runs the linter on files from a specific directory", lintHandler)
	cmd.AddFlag("path", "p", "DIR", "Path to .github directory", broccli.TypePathFile, broccli.IsDirectory|broccli.IsExistent|broccli.IsRequired)
	cmd.AddFlag("config", "c", "FILE", "Linter config with rules in YAML format", broccli.TypePathFile, broccli.IsRegularFile|broccli.IsExistent)
	cmd.AddFlag("loglevel", "l", "", "One of NONE,ERR,WARN,DEBUG", broccli.TypeString, 0)
	cmd.AddFlag("vars-file", "z", "", "Check if variable names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)
	cmd.AddFlag("secrets-file", "s", "", "Check if secret names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)

	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}

	os.Exit(cli.Run())
}

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func lintHandler(c *broccli.CLI) int {
	logLevel := loglevel.GetLogLevelFromString(c.Flag("loglevel"))
	varsFile := c.Flag("vars-file")
	secretsFile := c.Flag("secrets-file")

	lint := linter.Linter{
		LogLevel: logLevel,
	}
	dotGithub := dotgithub.DotGithub{
		LogLevel: logLevel,
	}

	err := dotGithub.ReadDir(c.Flag("path"))
	if err != nil {
		printErr(logLevel, err, "error initializing")
		return 20
	}

	if varsFile != "" {
		err = dotGithub.ReadVars(varsFile)
		if err != nil {
			printErr(logLevel, err, "error reading vars file")
			return 41
		}
	}
	if secretsFile != "" {
		err = dotGithub.ReadSecrets(secretsFile)
		if err != nil {
			printErr(logLevel, err, "error reading secrets file")
			return 42
		}
	}

	cfgFile, err := getConfigFilePath(c.Flag("config"), c.Flag("path"))
	if err != nil {
		printErr(logLevel, err, "error getting config file")
		return 21
	}

	cfg := linter.Config{
		LogLevel: logLevel,
	}
	if cfgFile != "" {
		err := cfg.ReadFile(cfgFile)
		if err != nil {
			printErr(logLevel, err, "error reading config file")
			return 31
		}
	} else {
		err := cfg.ReadDefaultFile()
		if err != nil {
			printErr(logLevel, err, "error reading default config file")
			return 32
		}
	}

	lint.Config = &cfg

	status, err := lint.Lint(&dotGithub)
	if err != nil {
		printErr(logLevel, err, "error linting")
		return 22
	}

	if status == linter.HasErrors {
		return 1
	}
	if status == linter.HasOnlyWarnings {
		return 2
	}

	return 0
}

func getConfigFilePath(filePath string, dotGitHubPath string) (string, error) {
	if filePath != "" {
		return filePath, nil
	}

	configInDotGithub := filepath.Join(dotGitHubPath, "dotgithub.yml")
	_, err := os.Stat(configInDotGithub)
	notFound := os.IsNotExist(err)
	if err != nil && !notFound {
		return "", fmt.Errorf("error getting os.Stat on dotgithub.yml inside .github path: %w", err)
	}
	if notFound {
		return "", nil
	}

	return configInDotGithub, nil
}

func printErr(logLevel int, err error, msg string) {
	if logLevel == loglevel.LogLevelNone {
		return
	}

	fmt.Fprintf(os.Stderr, "!!!:%s: %s\n", msg, err.Error())
}
