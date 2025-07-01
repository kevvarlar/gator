package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func write(c Config) error {
	jsonConfig, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshalize config into json: %w", err)
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}
	err = os.WriteFile(configFilePath, jsonConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}
	return nil
}