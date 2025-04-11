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
	Rules    []rule.Rule      `yaml:"-"`
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
	if len(cfg.Rules) > 0 {
		for _, r := range cfg.Rules {
			err := r.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
