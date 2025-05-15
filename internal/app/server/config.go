package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/i-Galts/go-server-project/internal/app/storage"
)

type ServerConfig struct {
	Port          string                `json:"port"`
	LogLevel      string                `json:"log_level"`
	CheckInterval string                `json:"check_interval"`
	Backends      []string              `json:"backends"`
	StorageConfig storage.StorageConfig `json:"storage"`
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

	fmt.Println("conf.StorageConfig.URL = ", conf.StorageConfig.URL)

	return conf, nil
}
