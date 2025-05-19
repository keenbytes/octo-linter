package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.pl/mikogs/broccli/v3"
	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/linter"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
)

func main() {
	cli := broccli.NewBroccli("octo-linter", "Validates GitHub Actions workflow and action YAML files", "m@gasior.dev")

	cmd := cli.Command("lint", "Runs the linter on files from a specific directory", lintHandler)
	cmd.Flag("path", "p", "DIR", "Path to .github directory", broccli.TypePathFile, broccli.IsDirectory|broccli.IsExistent|broccli.IsRequired)
	cmd.Flag("config", "c", "FILE", "Linter config with rules in YAML format", broccli.TypePathFile, broccli.IsRegularFile|broccli.IsExistent)
	cmd.Flag("loglevel", "l", "", "One of INFO,ERR,WARN,DEBUG", broccli.TypeString, 0)
	cmd.Flag("vars-file", "z", "", "Check if variable names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)
	cmd.Flag("secrets-file", "s", "", "Check if secret names exist in this file (one per line)", broccli.TypePathFile, broccli.IsExistent)

	_ = cli.Command("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}

	os.Exit(cli.Run(context.Background()))
}

func versionHandler(ctx context.Context, c *broccli.Broccli) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
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
