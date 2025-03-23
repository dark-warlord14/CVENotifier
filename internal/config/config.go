package config

import (
	"os"

	"github.com/dark-warlord14/CVENotifier/internal/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, &errors.ConfigError{Message: "Failed to read config file: " + err.Error()}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &errors.ConfigError{Message: "Failed to unmarshal config data: " + err.Error()}
	}

	return &cfg, nil
}

func GetConfigPath() string {
	configPath := "config.yaml"
	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		configPath = envConfigPath
	}
	return configPath
}
