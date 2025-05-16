// Package server provides structures and utilities for configuring and running the application server.
// This file focuses on loading server configuration from JSON files.
package server

import (
	"encoding/json"
	"os"
)

// holds all configurable parameters for the server
type ServerConfig struct {
	Port                  string   `json:"port"`
	LogLevel              string   `json:"log_level"`
	CheckInterval         string   `json:"check_interval"`
	Backends              []string `json:"backends"`
	RateLimiterCap        int      `json:"rl_capacity"`
	RateLimiterRefillRate int      `json:"rl_refillrate"`
}

// reads and parses a JSON configuration file into a ServerConfig struct
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
