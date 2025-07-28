package logger

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// New creates a new logger with the specified level
func New(logLevel LogLevel) log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Set the log level
	switch logLevel {
	case DebugLevel:
		logger = level.NewFilter(logger, level.AllowDebug())
	case InfoLevel:
		logger = level.NewFilter(logger, level.AllowInfo())
	case WarnLevel:
		logger = level.NewFilter(logger, level.AllowWarn())
	case ErrorLevel:
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger
}
