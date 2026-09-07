package utils

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/models"
)

// getXDGStateDir return dir $XDG_STATE_DIR
func getXDGStateDir() string {
	homeDir, err := models.GetHomeDir()
	if err != nil {
		slog.Error("while getting home dir for get xdg_state_dir", "err", err)
		log.Fatal(err)
	}

	return path.Join(homeDir, ".local", "state")
}

func getDefaultLogPath() string {
	return path.Join(getXDGStateDir(), models.AppName, "log")
}

var defaultLogPath = getDefaultLogPath()

// defaultGlobalLogger return the Logger with multiWrite
// the writer contains Stdout and [defaultLogPath]
// This default logger use [slog.TextHandler]
func defaultGlobalLogger() (logger *slog.Logger, logFile io.Closer) {
	err := EnsureFileExist(defaultLogPath)
	if err != nil {
		slog.Error(err.Error(), "func", "defaultGlobalLogger")
	}
	defaultLogFile, err := os.OpenFile(defaultLogPath, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		slog.Error("error when open default log file", "err", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, defaultLogFile)

	textHandler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{})
	return slog.New(textHandler), defaultLogFile
}

// SetGlobalLogger call [defaultGlobalLogger] as the slog default logger
// then return logFile which should be closed before program exit
// call this function on the root.go
// The you can just use slog for logging
func SetGlobalLogger() (logFile io.Closer) {
	logger, logFile := defaultGlobalLogger()

	slog.SetDefault(logger)

	return logFile
}
