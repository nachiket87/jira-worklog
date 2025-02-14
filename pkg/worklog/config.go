package worklog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Email   string `json:"email"`
	Token   string `json:"token"`
	BaseURL string `json:"base_url"`
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	// Creates ~/.config/jira-worklog/config.json
	configDir := filepath.Join(homeDir, ".config", "jira-worklog")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no configuration found. Please run 'worklog configure' first")
		}
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	// Set default base URL if not specified
	if config.BaseURL == "" {
		config.BaseURL = "https://onestepsoftware.atlassian.net/rest/api/3"
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	if config.Email == "" || config.Token == "" {
		return fmt.Errorf("email and token are required")
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	// Create file with user-only read/write permissions (0600)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
