package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// validateParams checks the validity of the creation type and name.
func ValidateParams(creationType, name string) error {
	if creationType == "" || (creationType != "file" && creationType != "dir") {
		return fmt.Errorf("invalid 'type' parameter. must be 'file' or 'dir'")
	}
	if name == "" {
		return fmt.Errorf("name parameter must be provided")
	}
	return nil
}

// resolvePath returns the cleaned and defaulted current path.
func ResolvePath(path string) string {
	if path == "" {
		return "." // Default to the current working directory
	}
	return filepath.Clean(path)
}

// isValidPath checks if the new path is within the base directory.
func IsValidPath(baseDir, newPath string) bool {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	absNewPath, err := filepath.Abs(newPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	return strings.HasPrefix(absNewPath, absBaseDir)
}

// createEntity creates a file or directory based on the type and path provided.
func CreateEntity(creationType, currentPath, name string) error {
	newPath := filepath.Join(currentPath, name)

	switch creationType {
	case "dir":
		if err := os.Mkdir(newPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create new directory: %w", err)
		}
	case "file":
		file, err := os.Create(newPath)
		if err != nil {
			return fmt.Errorf("failed to create new file: %w", err)
		}
		defer file.Close()
	default:
		return fmt.Errorf("unknown type %s", creationType)
	}
	return nil
}
