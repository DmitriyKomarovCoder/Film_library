package logger

import (
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logFilePath := "test.log"

	defer func() {
		_ = os.Remove(logFilePath)
	}()

	logger, err := NewLogger(logFilePath)
	if err != nil {
		t.Errorf("NewLogger() error = %v, want nil", err)
	}

	if logger == nil {
		t.Error("NewLogger() returned nil Logger, want non-nil Logger")
	}

	if logger.Logger.Out == nil {
		t.Error("NewLogger() did not set Logger's output")
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Errorf("NewLogger() did not create the log file %s", logFilePath)
	}

	file, err := os.OpenFile(logFilePath, os.O_WRONLY, 0666)
	if err != nil {
		t.Errorf("NewLogger() could not open the log file %s for writing: %v", logFilePath, err)
	}
	defer file.Close()

	logger.Logger.Info("Test log message")
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Failed to read the log file %s: %v", logFilePath, err)
	}

	if len(content) == 0 {
		t.Errorf("Log file %s is empty", logFilePath)
	}
}
