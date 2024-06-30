package services

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(isDebug bool) *logrus.Logger {
	logger := logrus.New()
	{
		if isDebug {
			logger.Level = logrus.DebugLevel
		} else {
			logger.Level = logrus.InfoLevel
		}
		logger.Formatter = &logrus.TextFormatter{}
		logger.Out = os.Stdout
	}
	return logger
}
