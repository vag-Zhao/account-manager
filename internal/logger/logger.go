package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level represents the log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a structured logger
type Logger struct {
	level      Level
	output     io.Writer
	mu         sync.Mutex
	fields     map[string]interface{}
	timeFormat string
}

// Global logger instance
var (
	defaultLogger *Logger
	once          sync.Once
)

// Initialize initializes the global logger
func Initialize(level Level, output io.Writer) {
	once.Do(func() {
		defaultLogger = &Logger{
			level:      level,
			output:     output,
			fields:     make(map[string]interface{}),
			timeFormat: "2006-01-02 15:04:05",
		}
	})
}

// Get returns the global logger instance
func Get() *Logger {
	if defaultLogger == nil {
		Initialize(InfoLevel, os.Stdout)
	}
	return defaultLogger
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:      l.level,
		output:     l.output,
		fields:     make(map[string]interface{}),
		timeFormat: l.timeFormat,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:      l.level,
		output:     l.output,
		fields:     make(map[string]interface{}),
		timeFormat: l.timeFormat,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// log writes a log entry
func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// Build log entry
	timestamp := time.Now().Format(l.timeFormat)
	entry := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level.String(), caller, msg)

	// Add fields if any
	if len(l.fields) > 0 {
		entry += " |"
		for k, v := range l.fields {
			entry += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	fmt.Fprintln(l.output, entry)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Package-level convenience functions

// Debug logs a debug message using the global logger
func Debug(msg string) {
	Get().Debug(msg)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(format string, args ...interface{}) {
	Get().Debugf(format, args...)
}

// Info logs an info message using the global logger
func Info(msg string) {
	Get().Info(msg)
}

// Infof logs a formatted info message using the global logger
func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

// Warn logs a warning message using the global logger
func Warn(msg string) {
	Get().Warn(msg)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

// Error logs an error message using the global logger
func Error(msg string) {
	Get().Error(msg)
}

// Errorf logs a formatted error message using the global logger
func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

// WithField adds a field to the global logger
func WithField(key string, value interface{}) *Logger {
	return Get().WithField(key, value)
}

// WithFields adds multiple fields to the global logger
func WithFields(fields map[string]interface{}) *Logger {
	return Get().WithFields(fields)
}

// SetLevel sets the log level for the global logger
func SetLevel(level Level) {
	Get().SetLevel(level)
}
