// Package logger is a nice package
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func SetupZerolog() *zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &logger
}
