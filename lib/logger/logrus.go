// Package logger is a nice package
package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func SetupLogrus() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
		PadLevelText: true,
	})
	return l
}
