package linter

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.pl/mikogs/octo-linter/pkg/linter/rule"
)

//go:embed dotgithub.yml
var defaultConfig []byte

type Config struct {
	Version        string                 `yaml:"version"`
	RulesConfig    map[string]interface{} `yaml:"rules"`
	ActionRules    []rule.ActionRule      `yaml:"-"`
	WorkflowRules  []rule.WorkflowRule    `yaml:"-"`
	DotGithubRules []rule.DotGithubRule   `yaml:"-"`
	Errors         map[string]string      `yaml:"errors"`
	LogLevel       int                    `yaml:"-"`
}

func (cfg *Config) ReadFile(p string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", p, err)
	}

	err = cfg.readBytesAndValidate(b)
	if err != nil {
		return fmt.Errorf("error reading and/or validating config file %s: %w", p, err)
	}

	return nil
}

func (cfg *Config) ReadDefaultFile() error {
	err := cfg.readBytesAndValidate(defaultConfig)
	if err != nil {
		return fmt.Errorf("error reading and/or validating default config file: %w", err)
	}

	return nil
}

func (cfg *Config) Validate() error {
	if len(cfg.ActionRules) > 0 {
		for _, r := range cfg.ActionRules {
			err := r.Validate()
			if err != nil {
				return err
			}
		}
	}
	if len(cfg.WorkflowRules) > 0 {
		for _, r := range cfg.WorkflowRules {
			err := r.Validate()
			if err != nil {
				return err
			}
		}
	}
	if len(cfg.DotGithubRules) > 0 {
		for _, r := range cfg.DotGithubRules {
			err := r.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
