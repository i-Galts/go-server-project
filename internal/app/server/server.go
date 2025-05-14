package server

import (
	"github.com/sirupsen/logrus"
)

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

func (s *ServerAPI) Launch() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("launching server...")

	return nil
}

func (s *ServerAPI) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}
