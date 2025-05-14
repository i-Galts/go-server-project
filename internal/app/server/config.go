package server

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Port          string   `json:"port"`
	LogLevel      string   `json:"log_level"`
	CheckInterval string   `json:"check_interval"`
	Backends      []string `json:"backends"`
}

func LoadConfig(file string) (ServerConfig, error) {
	var conf ServerConfig

	data, err := os.ReadFile(file)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
