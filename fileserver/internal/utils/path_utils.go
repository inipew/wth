package utils

import (
	"path/filepath"
	"strings"
)

// CleanPath returns a cleaned absolute path, ensuring it is within the base directory.
func CleanPath(basePath, inputPath string) (string, error) {
	// Clean and resolve the path
	absPath, err := filepath.Abs(filepath.Join(basePath, inputPath))
	if err != nil {
		return "", err
	}

	// Ensure the path is within the base directory
	baseAbsPath, err := filepath.Abs(basePath)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, baseAbsPath) {
		return "", nil // Access denied
	}

	return absPath, nil
}
