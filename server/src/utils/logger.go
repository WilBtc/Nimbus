// server/src/utils/logger.go

package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps the logrus.Logger with additional configurations
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new logger instance with optional log rotation and formatting settings.
func NewLogger(level, filePath string, rotate bool, rotationInterval int) *Logger {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warn("Invalid log level specified, defaulting to INFO")
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set formatter with full timestamps for detailed logs
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Configure log output
	if filePath != "" {
		// Log rotation settings using lumberjack
		logger.SetOutput(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    10,            // Max size in megabytes before rotation
			MaxBackups: 3,             // Max number of backups to retain
			MaxAge:     rotationInterval, // Max age in days before deletion
			Compress:   true,          // Compress old logs
		})
		logger.Info("Logging to file with rotation enabled")
	} else {
		// Default to stdout if no file path is provided
		logger.SetOutput(os.Stdout)
		logger.Info("Logging to stdout")
	}

	return &Logger{
		logger: logger,
	}
}

// Info logs informational messages
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Warn logs warning messages
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error logs error messages
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Debug logs debug messages, useful for development and troubleshooting
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Close performs any cleanup actions needed for the logger (currently handled by lumberjack).
func (l *Logger) Close() {
	// No explicit close needed, as lumberjack automatically handles file closures
}
