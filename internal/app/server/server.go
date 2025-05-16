// Package server provides the core server logic for the application.
// It includes configuration handling, logging setup, and the main server launch functionality.
package server

import (
	"github.com/sirupsen/logrus"
)

// represents the main server instance
type ServerAPI struct {
	config *ServerConfig
	logger *logrus.Logger
}

func NewServer(config *ServerConfig) *ServerAPI {
	return &ServerAPI{
		config: config,
		logger: logrus.New(),
	}
}

// starts the server and initializes all required components
func (s *ServerAPI) Launch() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	return nil
}

// sets up the logger's log level based on the server configuration
func (s *ServerAPI) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
