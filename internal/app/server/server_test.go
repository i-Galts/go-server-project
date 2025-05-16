package server_test

import (
	"os"
	"testing"

	"github.com/i-Galts/go-server-project/internal/app/server"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_ValidFile_ReturnsConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	configData := `{
        "port": "8080",
        "log_level": "info",
        "check_interval": "5s",
        "backends": ["http://backend1", "http://backend2"],
        "rl_capacity": 10,
        "rl_refillrate": 2
    }`

	err = os.WriteFile(tmpfile.Name(), []byte(configData), 0644)
	assert.NoError(t, err)

	conf, err := server.LoadConfig(tmpfile.Name())
	assert.NoError(t, err)

	assert.Equal(t, "8080", conf.Port)
	assert.Equal(t, "info", conf.LogLevel)
	assert.Equal(t, "5s", conf.CheckInterval)
	assert.Equal(t, []string{"http://backend1", "http://backend2"}, conf.Backends)
	assert.Equal(t, 10, conf.RateLimiterCap)
	assert.Equal(t, 2, conf.RateLimiterRefillRate)
}

func TestLoadConfig_InvalidJSON_ReturnsError(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	os.WriteFile(tmpfile.Name(), []byte("invalid json"), 0644)

	conf, err := server.LoadConfig(tmpfile.Name())
	assert.Error(t, err)
	assert.Empty(t, conf)
}

func TestLoadConfig_FileNotFound_ReturnsError(t *testing.T) {
	conf, err := server.LoadConfig("nonexistent.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
	assert.Empty(t, conf)
}

func TestLoadConfig_EmptyFields_DefaultValues(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	os.WriteFile(tmpfile.Name(), []byte("{}"), 0644)

	conf, err := server.LoadConfig(tmpfile.Name())
	assert.NoError(t, err)

	assert.Empty(t, conf.Port)
	assert.Empty(t, conf.LogLevel)
	assert.Empty(t, conf.CheckInterval)
	assert.Empty(t, conf.Backends)
	assert.Zero(t, conf.RateLimiterCap)
	assert.Zero(t, conf.RateLimiterRefillRate)
}
