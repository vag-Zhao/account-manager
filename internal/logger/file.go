package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter writes logs to a file with rotation support
type FileWriter struct {
	filename   string
	file       *os.File
	mu         sync.Mutex
	maxSize    int64 // Maximum size in bytes before rotation
	maxBackups int   // Maximum number of old log files to keep
}

// NewFileWriter creates a new file writer
func NewFileWriter(filename string, maxSize int64, maxBackups int) (*FileWriter, error) {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &FileWriter{
		filename:   filename,
		file:       file,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}, nil
}

// Write implements io.Writer
func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if rotation is needed
	if w.shouldRotate() {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

// shouldRotate checks if the log file should be rotated
func (w *FileWriter) shouldRotate() bool {
	if w.maxSize <= 0 {
		return false
	}

	info, err := w.file.Stat()
	if err != nil {
		return false
	}

	return info.Size() >= w.maxSize
}

// rotate rotates the log file
func (w *FileWriter) rotate() error {
	// Close current file
	if err := w.file.Close(); err != nil {
		return err
	}

	// Rotate old files
	for i := w.maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if i == w.maxBackups-1 {
				// Delete oldest file
				os.Remove(oldName)
			} else {
				// Rename file
				os.Rename(oldName, newName)
			}
		}
	}

	// Rename current file to backup
	if err := os.Rename(w.filename, w.backupName(0)); err != nil {
		return err
	}

	// Create new file
	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.file = file
	return nil
}

// backupName returns the backup filename for the given index
func (w *FileWriter) backupName(index int) string {
	if index == 0 {
		return w.filename + ".1"
	}
	return fmt.Sprintf("%s.%d", w.filename, index+1)
}

// Close closes the file writer
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

// MultiWriter writes to multiple writers
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter creates a new multi-writer
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

// Write implements io.Writer
func (w *MultiWriter) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		n, err = writer.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}

// InitializeWithFile initializes the logger with file output
func InitializeWithFile(level Level, filename string, maxSize int64, maxBackups int) error {
	fileWriter, err := NewFileWriter(filename, maxSize, maxBackups)
	if err != nil {
		return err
	}

	// Write to both file and stdout
	multiWriter := NewMultiWriter(fileWriter, os.Stdout)
	Initialize(level, multiWriter)

	return nil
}

// GetLogFilename returns the default log filename
func GetLogFilename() string {
	return filepath.Join("logs", fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
}
