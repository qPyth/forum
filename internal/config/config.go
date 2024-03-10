package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Database struct {
	PathToDb string `json:"path_to_db"`
}

func Load(pathToConfig string) (Config, error) {

	file, err := os.Open(pathToConfig)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config file not found: %w", err)
		}
	}

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}

	return config, nil
}
