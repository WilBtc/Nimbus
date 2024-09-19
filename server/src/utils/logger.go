// server/src/utils/logger.go

package utils

import (
	"io"
	"log"
	"os"
	"time"
)

// Logger struct provides a structured logging mechanism with multiple levels and file output options.
type Logger struct {
	logger     *log.Logger
	level      string
	logFile    *os.File
	filePath   string
	rotate     bool
	rotation   time.Duration
	nextRotate time.Time
}

// NewLogger creates a new Logger instance with the specified log level, file path, and rotation settings.
func NewLogger(level, filePath string, rotate bool, rotation time.Duration) *Logger {
	logger := &Logger{
		level:    level,
		filePath: filePath,
		rotate:   rotate,
		rotation: rotation,
	}

	logger.initLogger()
	return logger
}

// initLogger initializes the logger with the appropriate output settings, including file logging and rotation.
func (l *Logger) initLogger() {
	output := io.MultiWriter(os.Stdout)

	// Set up file logging if a file path is provided
	if l.filePath != "" {
		var err error
		l.logFile, err = os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		output = io.MultiWriter(output, l.logFile)
	}

	l.logger = log.New(output, "", log.LstdFlags)

	// Set up log rotation if enabled
	if l.rotate {
		l.nextRotate = time.Now().Add(l.rotation)
		go l.handleLogRotation()
	}
}

// Info logs informational messages if the log level is set appropriately.
func (l *Logger) Info(msg ...interface{}) {
	if l.level == "INFO" || l.level == "DEBUG" {
		l.logger.Println("[INFO]", msg...)
	}
}

// Error logs error messages regardless of the log level setting.
func (l *Logger) Error(msg ...interface{}) {
	l.logger.Println("[ERROR]", msg...)
}

// Debug logs debug messages if the log level is set to DEBUG.
func (l *Logger) Debug(msg ...interface{}) {
	if l.level == "DEBUG" {
		l.logger.Println("[DEBUG]", msg...)
	}
}

// handleLogRotation handles log file rotation based on the specified duration.
func (l *Logger) handleLogRotation() {
	for {
		time.Sleep(time.Minute) // Check every minute for log rotation
		if time.Now().After(l.nextRotate) {
			l.rotateLogFile()
			l.nextRotate = time.Now().Add(l.rotation)
		}
	}
}

// rotateLogFile performs the log file rotation by renaming the current log file and creating a new one.
func (l *Logger) rotateLogFile() {
	if l.logFile != nil {
		// Close current log file
		err := l.logFile.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close log file during rotation: %v", err)
		}

		// Rename the current log file with a timestamp
		rotatedName := l.filePath + "." + time.Now().Format("20060102_150405")
		err = os.Rename(l.filePath, rotatedName)
		if err != nil {
			l.logger.Printf("[ERROR] Failed to rotate log file: %v", err)
			return
		}

		// Open a new log file
		l.logFile, err = os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			l.logger.Fatalf("[ERROR] Failed to create new log file: %v", err)
		}

		// Update logger output to the new file
		output := io.MultiWriter(os.Stdout, l.logFile)
		l.logger.SetOutput(output)
		l.logger.Printf("[INFO] Log file rotated: %s", rotatedName)
	}
}

// Close gracefully shuts down the logger, closing any open files.
func (l *Logger) Close() {
	if l.logFile != nil {
		err := l.logFile.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close log file: %v", err)
		}
	}
}
