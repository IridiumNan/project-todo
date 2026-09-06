package config

import (
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/models"
)

const maxSearchDepth = 4

// searchDataDir search the dir .project-todo dir
// search for [maxSearchDepth]
// if not found, return the [os.ErrNotExist]
// Then you should find data or configuration file on another path
func searchDataDir() (string, error) {
	searchCount := 0

	currPath, err := os.Getwd()
	if err != nil {
		slog.Error("while getting current dir path", "err", err)
		currPath = "."
	}

	for searchCount <= maxSearchDepth {
		if containsDataDir(currPath) {
			return currPath, nil
		}

		slog.Debug("data dir not found for current path, checking next", "current_path", currPath, "target", models.DataDirName)

		currPath = path.Dir(currPath)

		searchCount++
	}
	return emptyPath, os.ErrNotExist
}

// containsDataDir is used for check if .AppName dir for data and configuration storage exist
// if not exist, it will return false
// if encounter err when readDir, it will write a log into log file then return false
func containsDataDir(checkedPath string) bool {
	entries, err := os.ReadDir(checkedPath)
	if err != nil {
		slog.Error("while reading dir", "dir_path", checkedPath, "err", err)
		return false
	}

	for idx := range entries {
		if entries[idx].IsDir() && entries[idx].Name() == models.DataDirName {
			return true
		}
	}

	return false
}
