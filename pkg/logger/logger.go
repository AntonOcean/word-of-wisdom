package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	log  *logrus.Logger
	once sync.Once
)

// Init initializes the logger once
func Init() {
	once.Do(func() {
		log = logrus.New()
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(logrus.DebugLevel)
	})
}

// GetLogger returns a singleton logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}
