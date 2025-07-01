package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)
const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory: %w", err)
	}
	return home + "/gator/" + configFileName, nil
}

func Read() (Config, error) {
	config := Config{}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("failed to get config file path: %w", err)
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to decode config file: %w", err)
	}
	fmt.Println(config)
	return config, nil
}