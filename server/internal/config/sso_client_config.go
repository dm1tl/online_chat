package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type SSOConfig struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount uint          `yaml:"retriescount"`
}

func NewSSOConfig() (*SSOConfig, error) {
	var config SSOConfig
	filepath := os.Getenv("SSO_CONFIG_PATH")
	if filepath == "" {
		return nil, fmt.Errorf("SSO_CONFIG_PATH is empty")
	}
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config file %w", err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("couldn't parse config file into model %w", err)
	}
	return &config, nil
}
