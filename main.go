package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/keenbytes/broccli/v3"
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
	"github.com/keenbytes/octo-linter/pkg/linter"
	"github.com/keenbytes/octo-linter/pkg/loglevel"
)

const configFileName = "dotgithub.yml"

func main() {
	cli := broccli.NewBroccli("octo-linter", "Validates GitHub Actions workflow and action YAML files", "m@gasior.dev")

	cmdLint := cli.Command("lint", "Runs the linter on files from a specific directory", lintHandler)
	cmdLint.Flag("path", "p", "DIR", "Path to .github directory", broccli.TypePathFile, broccli.IsDirectory|broccli.IsExistent|broccli.IsRequired)
	cmdLint.Flag("config", "c", "FILE", "Linter config with rules in YAML format", broccli.TypePathFile, broccli.IsRegularFile|broccli.IsExistent)
	cmdLint.Flag("", "v", "", "Prints more information", broccli.TypeBool, 0)
	cmdLint.Flag("", "vv", "", "Prints debug information", broccli.TypeBool, 0)
	cmdLint.Flag("quiet", "q", "", "No output", broccli.TypeBool, 0)
	cmdLint.Flag("vars-file", "z", "", "Check if variable names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)
	cmdLint.Flag("secrets-file", "s", "", "Check if secret names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)

	cmdInit := cli.Command("init", "Create sample dotgithub.yml config file", initHandler)
	cmdInit.Flag("destination", "d", "FILE", "Destination filename to write to", broccli.TypePathFile, broccli.IsNotExistent)

	_ = cli.Command("version", "Prints version", versionHandler)

	os.Exit(cli.Run(context.Background()))
}

func versionHandler(ctx context.Context, c *broccli.Broccli) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func initHandler(ctx context.Context, c *broccli.Broccli) int {
	path := c.Flag("destination")
	if path == "" {
		_, err := os.Stat(configFileName)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("error checking if destination path exists: %s", err.Error()))
				return 50
			}
		} else {
			slog.Error(fmt.Sprintf("file %s already exists, remove it or use --destination flag to change the destination", configFileName))
			return 51
		}
		path = configFileName
	}

	err := os.WriteFile(path, linter.GetDefaultConfig(), 0644)
	if err != nil {
		slog.Error(fmt.Sprintf("error checking if destination path exists: %s", err.Error()))
		return 52
	}

	slog.Info(fmt.Sprintf("Sample configuration file %s has been created. Run 'lint' command with '-c' flag or put the file in the .github directory.", path))
	return 0
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
		return 20
	}

	if varsFile != "" {
		err = dotGithub.ReadVars(varsFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading vars file: %s", err.Error()))
			return 41
		}
	}
	if secretsFile != "" {
		err = dotGithub.ReadSecrets(secretsFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading secrets file: %s", err.Error()))
			return 42
		}
	}

	cfgFile, err := getConfigFilePath(c.Flag("config"), c.Flag("path"))
	if err != nil {
		slog.Error(fmt.Sprintf("error getting config file: %s", err.Error()))
		return 21
	}

	cfg := linter.Config{}
	if cfgFile != "" {
		err := cfg.ReadFile(cfgFile)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading config file: %s", err.Error()))
			return 31
		}
	} else {
		err := cfg.ReadDefaultFile()
		if err != nil {
			slog.Error(fmt.Sprintf("error reading default config file: %s", err.Error()))
			return 32
		}
	}

	lint.Config = &cfg

	status, err := lint.Lint(&dotGithub)
	if err != nil {
		slog.Error(fmt.Sprintf("error linting: %s", err.Error()))
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
