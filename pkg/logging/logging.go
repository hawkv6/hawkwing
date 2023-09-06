package logging

import "github.com/sirupsen/logrus"

var DefaultLogger = InitializeDefaultLogger()

func InitializeDefaultLogger() (logger *logrus.Logger) {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	return
}
