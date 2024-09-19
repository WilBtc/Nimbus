// server/src/utils/logger.go

package utils

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps the logrus.Logger with additional configurations
type Logger struct {
	logger *logrus.Logger
}

func NewLogger(level, filePath string, rotate bool, rotationInterval int) *Logger {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warn("Invalid log level specified, defaulting to INFO")
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set output
	if filePath != "" {
		// Configure log rotation
		logger.SetOutput(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    10, // Max size in megabytes
			MaxBackups: 3,
			MaxAge:     rotationInterval, // Max age in days
			Compress:   true,
		})
	} else {
		logger.SetOutput(os.Stdout)
	}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Close() {
	// No action needed as lumberjack handles file closing
}
