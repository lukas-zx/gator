package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func getConfigPath() string {
	userHome, err := os.UserHomeDir()
	panicOnError(err)
	return path.Join(userHome, configFileName)
}

func Read() *Config {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	panicOnError(err)

	config := &Config{}
	err = json.Unmarshal(data, config)
	panicOnError(err)

	return config
}

func (config *Config) SetUser(username string) {
	config.CurrentUserName = username
	err := writeConfig(*config)
	panicOnError(err)
}

func writeConfig(config Config) error {
	json, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	path := getConfigPath()
	err = os.WriteFile(path, json, 0644)
	if err != nil {
		return fmt.Errorf("could write config: %w", err)
	}

	return nil
}
