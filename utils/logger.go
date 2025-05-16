package utils

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is the global logger instance
var Logger zerolog.Logger

func init() {
	// Initialize logger
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	Logger = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}

// SetLogLevel sets the logging level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		Logger = Logger.Level(zerolog.DebugLevel)
	case "info":
		Logger = Logger.Level(zerolog.InfoLevel)
	case "warn":
		Logger = Logger.Level(zerolog.WarnLevel)
	case "error":
		Logger = Logger.Level(zerolog.ErrorLevel)
	default:
		Logger = Logger.Level(zerolog.InfoLevel)
	}
}
