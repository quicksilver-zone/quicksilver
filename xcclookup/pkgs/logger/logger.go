package logger

import (
	"context"
	"os"

	kitlog "github.com/go-kit/log"
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

// Logger wraps the go-kit logger with additional convenience methods
type Logger struct {
	logger kitlog.Logger
}

// New creates a new logger with the specified level
func New(logLevel LogLevel) *Logger {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

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

	return &Logger{logger: logger}
}

// With returns a new logger with the given key-value pairs added to the context
func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{logger: kitlog.With(l.logger, keyvals...)}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	level.Debug(l.logger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

// Info logs an info message
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	level.Info(l.logger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	level.Warn(l.logger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	level.Error(l.logger).Log(append([]interface{}{"msg", msg}, keyvals...)...)
}

// GetKitLogger returns the underlying go-kit logger for advanced usage
func (l *Logger) GetKitLogger() kitlog.Logger {
	return l.logger
}

// Context key for logger
type loggerKey struct{}

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext retrieves a logger from the context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return logger
	}
	// Return a default logger if none is found in context
	return New(InfoLevel)
}
