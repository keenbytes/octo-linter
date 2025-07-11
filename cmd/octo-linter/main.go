package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/keenbytes/broccli/v3"
	"github.com/keenbytes/octo-linter/v2/internal/linter"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/v2/pkg/loglevel"
)

//go:generate go run ../../gen.go ../../

const configFileName = "dotgithub.yml"

// exit codes
const (
	ExitOK                       = 0
	ExitLintErrors               = 1
	ExitLintOnlyWarnings         = 2
	ExitErrLinting               = 10
	ExitErrReadingDotGithubDir   = 20
	ExitErrGettingCfgFile        = 30
	ExitErrReadingCfgFile        = 31
	ExitErrReadingDefaultCfgFile = 32
	ExitErrReadingVarsFile       = 41
	ExitErrReadingSecretsFile    = 42
	ExitErrCheckingDstPath       = 50
	ExitDstFileIsDir             = 51
	ExitErrWritingCfg            = 52
)

func main() {
	cli := broccli.NewBroccli("octo-linter", "Validates GitHub Actions workflow and action YAML files", "m@gasior.dev")

	cmdLint := cli.Command("lint", "Runs the linter on files from a specific directory", lintHandler)
	cmdLint.Flag("path", "p", "DIR", "Path to .github directory", broccli.TypePathFile, broccli.IsDirectory|broccli.IsExistent|broccli.IsRequired)
	cmdLint.Flag("config", "c", "FILE", "Linter config with rules in YAML format", broccli.TypePathFile, broccli.IsRegularFile|broccli.IsExistent)
	cmdLint.Flag("loglevel", "l", "", "One of INFO,ERR,WARN,DEBUG", broccli.TypeString, 0)
	cmdLint.Flag("vars-file", "z", "", "Check if variable names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)
	cmdLint.Flag("secrets-file", "s", "", "Check if secret names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)
	cmdLint.Flag("output", "o", "DIR", "Path to where summary markdown gets generated", broccli.TypePathFile, broccli.IsDirectory|broccli.IsExistent)
	cmdLint.Flag("output-errors", "u", "INT", "Limit numbers of errors shown in the markdown output file", broccli.TypeInt, 0)

	cmdInit := cli.Command("init", "Create sample dotgithub.yml config file", initHandler)
	cmdInit.Flag("destination", "d", "FILE", "Destination filename to write to", broccli.TypePathFile, broccli.IsNotExistent)

	_ = cli.Command("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}

	os.Exit(cli.Run(context.Background()))
}

func versionHandler(ctx context.Context, c *broccli.Broccli) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return ExitOK
}

func initHandler(ctx context.Context, c *broccli.Broccli) int {
	path := c.Flag("destination")
	if path == "" {
		fileInfo, err := os.Stat(configFileName)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("error checking if destination path exists: %s", err.Error()))
				return ExitErrCheckingDstPath
			}
		} else {
			if fileInfo.IsDir() {
				slog.Error(fmt.Sprintf("file %s already exists and it is a directory, remove it first or use --destination flag to change the destination", configFileName))
				return ExitDstFileIsDir
			}
		}
		path = configFileName
	}

	err := os.WriteFile(path, linter.GetDefaultConfig(), 0644)
	if err != nil {
		slog.Error(fmt.Sprintf("error writing default config: %s", err.Error()))
		return ExitErrWritingCfg
	}

	slog.Info(fmt.Sprintf("Sample configuration file %s has been created. Run 'lint' command with '-c' flag or put the file in the .github directory.", path))
	return ExitOK
}

func lintHandler(ctx context.Context, c *broccli.Broccli) int {
	logLevel := loglevel.GetLogLevelFromString(c.Flag("loglevel"))
	varsFile := c.Flag("vars-file")
	secretsFile := c.Flag("secrets-file")

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	lint := linter.Linter{}
	dotGithub := dotgithub.DotGithub{}

	err := dotGithub.ReadDir(c.Flag("path"))
	if err != nil {
		slog.Error(fmt.Sprintf("error initializing: %s", err.Error()))
		return ExitErrReadingDotGithubDir
	}

	if varsFile != "" {
		err = dotGithub.ReadVars(varsFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading vars file: %s", err.Error()))
			return ExitErrReadingVarsFile
		}
	}
	if secretsFile != "" {
		err = dotGithub.ReadSecrets(secretsFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading secrets file: %s", err.Error()))
			return ExitErrReadingSecretsFile
		}
	}

	cfgFile, err := getConfigFilePath(c.Flag("config"), c.Flag("path"))
	if err != nil {
		slog.Error(fmt.Sprintf("error getting config file: %s", err.Error()))
		return ExitErrGettingCfgFile
	}

	cfg := linter.Config{}
	if cfgFile != "" {
		err := cfg.ReadFile(cfgFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading config file: %s", err.Error()))
			return ExitErrReadingCfgFile
		}
	} else {
		err := cfg.ReadDefaultFile()
		if err != nil {
			slog.Error(fmt.Sprintf("error reading default config file: %s", err.Error()))
			return ExitErrReadingDefaultCfgFile
		}
	}

	lint.Config = &cfg

	outputLimit := 0
	if c.Flag("output-errors") != "" {
		// flag is already validated by the cli
		outputLimit, _ = strconv.Atoi(c.Flag("output-errors"))
	}

	status, err := lint.Lint(&dotGithub, c.Flag("output"), outputLimit)
	if err != nil {
		slog.Error(fmt.Sprintf("error linting: %s", err.Error()))
		return ExitErrLinting
	}

	if status == linter.HasErrors {
		return ExitLintErrors
	}
	if status == linter.HasOnlyWarnings {
		return ExitLintOnlyWarnings
	}

	return 0
}

func getConfigFilePath(filePath string, dotGitHubPath string) (string, error) {
	if filePath != "" {
		return filePath, nil
	}

	configInDotGithub := filepath.Join(dotGitHubPath, configFileName)
	_, err := os.Stat(configInDotGithub)
	notFound := os.IsNotExist(err)
	if err != nil && !notFound {
		return "", fmt.Errorf("error getting os.Stat on %s inside .github path: %w", configFileName, err)
	}
	if notFound {
		return "", nil
	}

	return configInDotGithub, nil
}
