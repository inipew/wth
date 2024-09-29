// file: fileutils.go
package fileutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sbx/internal/logger"
	"sync"
)

var (
	ErrNotDirectory = errors.New("not a directory")
	ErrPathNotExist = errors.New("path does not exist")
)

// Exists checks if the given path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDir creates a directory with the specified permissions.
func CreateDir(dir string, perm os.FileMode) error {
	if err := os.MkdirAll(dir, perm); err != nil {
		logger.GetLogger().Error().Err(err).Str("directory", dir).Msg("Failed to create directory")
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	logger.GetLogger().Info().Msgf("Directory created: %s", dir)
	return nil
}

// GetPermissions returns the permissions of the specified path in numeric format.
func GetPermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get permissions: %w", err)
	}
	return fmt.Sprintf("%04o", info.Mode().Perm()), nil
}

// SetPermissions sets the permissions of the specified path.
func SetPermissions(path string, perm os.FileMode) error {
	if err := os.Chmod(path, perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}
	logger.GetLogger().Info().Str("permission", fmt.Sprintf("%04o", perm)).Str("filepath", path).Msg("Successfully set permission")
	return nil
}

// Remove removes a file or directory.
func Remove(path string) error {
	if !Exists(path) {
		return ErrPathNotExist
	}
	logger.GetLogger().Info().Str("file",path).Msg("File removed.")
	return os.RemoveAll(path)
}

// Rename renames a file or directory.
func Rename(oldPath, newPath string) error {
	if !Exists(oldPath) {
		return ErrPathNotExist
	}
	return os.Rename(oldPath, newPath)
}

// ListFiles returns a list of files in the specified directory.
func ListFiles(dir string) ([]string, error) {
	if !Exists(dir) {
		return nil, ErrPathNotExist
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames, nil
}

// GetFileSize returns the size of the specified file.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}
	return info.Size(), nil
}

// IsDirectory checks if the given path is a directory.
func IsDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if directory: %w", err)
	}
	return info.IsDir(), nil
}

// IsFile checks if the given path is a file.
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if file: %w", err)
	}
	return !info.IsDir(), nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// MoveFile moves a file from src to dst.
func MoveFile(src, dst string) error {
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

// CopyDir copies a directory recursively.
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(fds))

	for _, fd := range fds {
		wg.Add(1)
		go func(fd os.DirEntry) {
			defer wg.Done()

			srcfp := filepath.Join(src, fd.Name())
			dstfp := filepath.Join(dst, fd.Name())

			if fd.IsDir() {
				if err := CopyDir(srcfp, dstfp); err != nil {
					errChan <- err
				}
			} else {
				if err := CopyFile(srcfp, dstfp); err != nil {
					errChan <- err
				}
			}
		}(fd)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
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
	isDir, err := IsDirectory(path)
	if err != nil {
		return err
	}
	if !isDir {
		return ErrNotDirectory
	}
	if !Exists(path) {
		return ErrPathNotExist
	}
	return os.Remove(path)
}