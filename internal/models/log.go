package models

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/utils"
)

// getXDGStateDir return dir $XDG_STATE_DIR
func getXDGStateDir() string {
	homeDir, err := GetHomeDir()
	if err != nil {
		slog.Error("while getting home dir for get xdg_state_dir", "err", err)
		log.Fatal(err)
	}

	return path.Join(homeDir, ".local", "state")
}

func getDefaultLogPath() string {
	return getXDGStateDir() + AppName
}

var defaultLogPath = getDefaultLogPath()

var GlobalLog *slog.Logger

// defaultGlobalLogger return the Logger with multiWrite
// the writer contains Stdout and [defaultLogPath]
// This default logger use [slog.TextHandler]
func defaultGlobalLogger() (logger *slog.Logger, logFile *os.File) {
	err := utils.EnsureFileExist(defaultLogPath)
	if err != nil {
		slog.Error(err.Error(), "func", "defaultGlobalLogger")
	}
	defaultLogFile, err := os.OpenFile(defaultLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)

	multiWriter := io.MultiWriter(os.Stdout, defaultLogFile)

	textHandler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{})
	return slog.New(textHandler), defaultLogFile
}

// SetGlobalLogger call [defaultGlobalLogger] as the slog default logger
// then return logFile which should be closed before program exit
// call this function on the root.go
// The you can just use slog for logging
func SetGlobalLogger() (logFile *os.File) {
	logger, logFile := defaultGlobalLogger()

	slog.SetDefault(logger)

	return logFile
}
