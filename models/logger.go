package models

import (
	"fmt"
	"os"
	"time"
)

// EggLog struct holds the file and log writers
type EggLog struct {
	file *os.File
}

// NewLogger creates a new logger instance
func NewLogger(filepath string) (*EggLog, error) {
	// Open the file in append mode (a+)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &EggLog{
		file: f,
	}, nil
}

// Close closes the log file
func (l *EggLog) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Info logs information messages
func (l *EggLog) Info(message string, args ...any) error {
	// Format the message using fmt.Sprintf and include all args
	message = fmt.Sprintf(message, args...)
	formatted := fmt.Sprintf("%s INFO: %s", time.Now().Format("2006-01-02 15:04:05.000"), message)

	// Write the formatted message to the log file with all args included
	_, err := l.file.WriteString(formatted + "\n")
	if err != nil {
		return err
	}
	return l.file.Sync()
}

// Error logs error messages with variable arguments
func (l *EggLog) Error(message string, args ...any) error {
	// Format the message using fmt.Sprintf and include all args
	message = fmt.Sprintf(message, args...)
	formatted := fmt.Sprintf("%s ERROR: %s", time.Now().Format("2006-01-02 15:04:05.000"), message)

	// Write the formatted message to the log file with all args included
	_, err := l.file.WriteString(formatted + "\n")
	if err != nil {
		return err
	}
	return l.file.Sync()
}
