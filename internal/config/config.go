// Package config used for pulling in users config options from yaml file
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Models          map[string]ProviderConfig `yaml:"providers"`
	DefaultProvider string                    `yaml:"default-provider"`
	DefaultModel    string                    `yaml:"default-model"`
}

type ProviderConfig struct {
	APIKey string   `yaml:"api_key"`
	Models []string `yaml:"models"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot load home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "ai-cli", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	return &cfg, nil
}
