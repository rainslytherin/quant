package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var configFiles = []string{"config_private.json", "config.json"}

var globalConfig *Config = nil

func GetGlobalConfig() *Config {
	return globalConfig
}

func InitGlobalConfig(dir string) {
	cfg, err := ReadConfig(dir)
	if err != nil {
		panic(err)
	}
	globalConfig = cfg
}

// read config from config_private.json and config.json
// and return config object

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
}

func ReadConfig(dir string) (*Config, error) {
	for _, configFile := range configFiles {
		config, err := readConfigFile(filepath.Join(dir, configFile))
		if err == nil {
			fmt.Printf("Using configuration from %s\n", configFile)
			return config, nil
		}
	}

	return &Config{}, errors.New("No valid configuration file found. Please create config.json or config_private.json")
}

func readConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
