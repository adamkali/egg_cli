/*
Copyright © 2025 Adam Kalinowski <adam.kalilarosa@proton.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
