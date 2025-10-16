// Package logger provides centralized logging functionality for the game
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger provides different log levels and file output
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	file        *os.File
}

// LogLevel represents the severity of log messages
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	WARM
)

var (
	// Global logger instance
	gameLogger *Logger
	// Verbose controls whether detailed/frequent logs are output
	Verbose bool = false
)

// Init initializes the logger with a file output
func Init() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Clean up any existing log files
	matches, _ := filepath.Glob(filepath.Join(logsDir, "myrpg_*.log"))
	for _, match := range matches {
		os.Remove(match)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logsDir, fmt.Sprintf("myrpg_%s.log", timestamp))

	// Create fresh log file (O_TRUNC ensures it starts clean)
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	gameLogger = &Logger{
		file:        file,
		debugLogger: log.New(file, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		infoLogger:  log.New(file, "[INFO]  ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		errorLogger: log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		warnLogger:  log.New(file, "[WARN]  ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}

	// Also log to stdout for immediate feedback
	Info("Logger initialized, log file: %s", logFileName)
	fmt.Printf("Logger initialized, log file: %s\n", logFileName)

	return nil
}

// Close closes the log file
func Close() error {
	if gameLogger != nil && gameLogger.file != nil {
		return gameLogger.file.Close()
	}
	return nil
}

// Debug logs debug messages (most verbose)
func Debug(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.debugLogger.Printf(format, args...)
	}
}

// Info logs informational messages
func Info(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.infoLogger.Printf(format, args...)
	}
}

// Warn logs warning messages
func Warn(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.warnLogger.Printf(format, args...)
	}
}

// Error logs error messages
func Error(format string, args ...interface{}) {
	if gameLogger != nil {
		gameLogger.errorLogger.Printf(format, args...)
	}
}

// Combat logs combat-specific debug messages with special prefix
func Combat(format string, args ...interface{}) {
	if gameLogger != nil {
		message := fmt.Sprintf("[COMBAT] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// UI logs UI-specific debug messages with special prefix
func UI(format string, args ...interface{}) {
	if gameLogger != nil {
		message := fmt.Sprintf("[UI] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// Turn logs turn management debug messages with special prefix
func Turn(format string, args ...interface{}) {
	if gameLogger != nil {
		message := fmt.Sprintf("[TURN] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// Action logs action execution debug messages with special prefix
func Action(format string, args ...interface{}) {
	if gameLogger != nil {
		message := fmt.Sprintf("[ACTION] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// VerboseDebug logs debug messages only if Verbose mode is enabled
func VerboseDebug(format string, args ...interface{}) {
	if gameLogger != nil && Verbose {
		gameLogger.debugLogger.Printf(format, args...)
	}
}

// VerboseCombat logs combat messages only if Verbose mode is enabled
func VerboseCombat(format string, args ...interface{}) {
	if gameLogger != nil && Verbose {
		message := fmt.Sprintf("[COMBAT] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// VerboseTurn logs turn messages only if Verbose mode is enabled
func VerboseTurn(format string, args ...interface{}) {
	if gameLogger != nil && Verbose {
		message := fmt.Sprintf("[TURN] "+format, args...)
		gameLogger.debugLogger.Print(message)
	}
}

// SetVerbose enables or disables verbose logging
func SetVerbose(enabled bool) {
	Verbose = enabled
}

// GetLogFileName returns the current log file name (for user reference)
func GetLogFileName() string {
	if gameLogger != nil && gameLogger.file != nil {
		return gameLogger.file.Name()
	}
	return ""
}
