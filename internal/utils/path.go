package utils

import (
	"fmt"
	"os"
	"path"
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
