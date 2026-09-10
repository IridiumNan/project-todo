package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/IridiumNan/project-todo/internal/models"
)

// EnsureFileExist utils function which ensure file exist
// It will automatically create dir and file to ensure that
func EnsureFileExist(filePath string) error {
	// Check if file exist
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	dirPath := path.Dir(filePath)

	err := os.MkdirAll(dirPath, 0o755)
	if err != nil {
		return fmt.Errorf("while mkdir %s, err: %w", dirPath, err)
	}

	_, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("while create file: %s, err: %w", filePath, err)
	}

	return nil
}

const (
	maxSearchDepth = 4
	emptyPath      = ""
)

// SearchDir search for target dir
// begin with startDir
// if you want to begin with current dir, use [os.Getwd]
// if targetDir not found, it will return [os.ErrNotExist]
func SearchDir(targetDir string, startDir string) (foundPath string, err error) {
	searchCount := 0

	rootDir := "/"

	currPath := startDir

	for searchCount <= maxSearchDepth && currPath != rootDir {
		if isContainsDir(currPath, targetDir) {
			return currPath, nil
		}

		slog.Debug("data dir not found for current path, checking next", "current_path", currPath, "target", models.DataDirName)

		currPath = path.Dir(currPath)

		searchCount++
	}
	return emptyPath, os.ErrNotExist
}

// isContainsDir check if checkedPath (as dir) contains the targetDir
// It just called by [SearchDir]
func isContainsDir(checkedPath string, targetDir string) bool {
	entries, err := os.ReadDir(checkedPath)
	if err != nil {
		slog.Error("while reading dir", "dir_path", checkedPath, "err", err)
		return false
	}

	for idx := range entries {
		if entries[idx].IsDir() && entries[idx].Name() == targetDir {
			return true
		}
	}

	return false
}

// ParseCtxPath return the id and title based on context file path
func ParseCtxPath(ctxPath string) (id string, title string) {
	ctxFileName := path.Base(ctxPath)
	parts := strings.SplitN(ctxFileName, "-", 2)

	return parts[0], parts[1]
}
