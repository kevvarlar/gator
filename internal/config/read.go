package config

import (
	"fmt"
	"os"
)
const configFileName = ".gatorconfig.json"
func Read() (Config, error) {
	config := Config{}
	home, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("unable to get home directory")
	}
	fmt.Println(home)
	return config, nil
}