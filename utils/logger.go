package utils

import (
	"github.com/sirupsen/logrus"
)

func SetupLogger(config *AppConfig) *logrus.Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Errorf("Invalid log level: %v, using default log level: info", err)
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	return log
}