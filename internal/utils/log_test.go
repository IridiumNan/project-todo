package utils

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"testing"
)

func clearOldLog(t *testing.T) {
	if err := os.RemoveAll(path.Dir(defaultLogPath)); err != nil {
		slog.Error("when removing old log file", "err", err)
		t.Errorf("error when removing old log file %s", err.Error())
	}
}

func TestLoggerCreate(t *testing.T) {
	clearOldLog(t)

	logFile := SetGlobalLogger()

	// Write some log into log file then check it's content

	slog.Info("test writing log", "err", errors.New("test writing error log"))

	if _, err := os.Stat(defaultLogPath); os.IsNotExist(err) {
		t.Errorf("log file not found, err: %s", err.Error())
	}

	defer logFile.Close()

	clearOldLog(t)
}
