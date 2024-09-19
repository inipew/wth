// file: fileutils.go
package fileutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// Exists checks if the given path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDir creates a directory with the specified permissions.
func CreateDir(dir string, perm os.FileMode) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, perm); err != nil {
			log.Error().Err(err).Str("directory", dir).Msg("Failed to create directory")
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Info().Msgf("Directory created: %s", dir)
	}
	return nil
}

// GetPermissions returns the permissions of the specified path in numeric format.
func GetPermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04o", info.Mode().Perm()), nil
}

// SetPermissions sets the permissions of the specified path.
func SetPermissions(path string, perm os.FileMode) error {
	log.Info().Str("permission", fmt.Sprintf("%o", perm)).Str("filepath", path).Msg("Successfully set permission")
	return os.Chmod(path, perm)
}

// Remove removes a file or directory.
func Remove(path string) error {
	if !Exists(path) {
		return os.ErrNotExist
	}
	return os.RemoveAll(path)
}

// Rename renames a file or directory.
func Rename(oldPath, newPath string) error {
	if !Exists(oldPath) {
		return os.ErrNotExist
	}
	return os.Rename(oldPath, newPath)
}

// ListFiles returns a list of files in the specified directory.
func ListFiles(dir string) ([]string, error) {
	if !Exists(dir) {
		return nil, os.ErrNotExist
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

// GetFileSize returns the size of the specified file.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// IsDirectory checks if the given path is a directory.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsFile checks if the given path is a file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	if !Exists(src) {
		return os.ErrNotExist
	}

	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

// MoveFile moves a file from src to dst.
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// CopyDir copies a directory recursively.
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return CopyFile(path, destPath)
	})
}

// MoveDir moves a directory recursively.
func MoveDir(src string, dst string) error {
	if err := CopyDir(src, dst); err != nil {
		return err
	}
	return Remove(src)
}

// RemoveEmptyDir removes a directory only if it is empty.
func RemoveEmptyDir(path string) error {
	if !IsDirectory(path) {
		return errors.New("not a directory")
	}
	if !Exists(path) {
		return os.ErrNotExist
	}
	return os.Remove(path)
}
