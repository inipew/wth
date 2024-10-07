package file

import (
	"errors"
	"files/internal/config"
	"files/internal/models"
	"files/internal/utils/helper"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	defaultFilePermissions = 0755
	timeFormat             = "2006-01-02 15:04:05"
)

// FileManager handles file operations using the application config
type FileManager struct {
	Config *config.Config
}

// NewFileManager creates a new FileManager instance
func NewFileManager(cfg *config.Config) *FileManager {
	return &FileManager{Config: cfg}
}

// GetDirectoryPath extracts and validates the directory path from the request
func (fm *FileManager) GetDirectoryPath(c *fiber.Ctx) (string, error) {
	dir := c.Query("path")
	if dir == "" {
		return fm.Config.Files.StorageDir, nil
	}

	decodedDir, err := url.QueryUnescape(dir)
	if err != nil {
		return "", fmt.Errorf("failed to unescape directory path: %w", err)
	}
	cleanedDir := filepath.Clean(decodedDir)
	absPath, err := filepath.Abs(cleanedDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	if !strings.HasPrefix(absPath, fm.Config.Files.StorageDir) {
		return "", errors.New("access to directory outside of storage area is forbidden")
	}

	return absPath, nil
}

// PrepareFileInfo prepares FileInfo slice from directory entries
func (fm *FileManager) PrepareFileInfo(files []os.DirEntry, dirPath string) ([]models.FileInfo, error) {
	var fileInfos []models.FileInfo

	for _, file := range files {
		if !fm.Config.Files.ShowHiddenFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("Error getting file info")
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		if info.Size() > fm.Config.Files.MaxFileSize {
			log.Warn().Str("file", file.Name()).Msg("File exceeds maximum size limit")
			continue
		}

		owner, group, err := fm.getFileOwnerGroup(filePath)
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("Error getting file owner and group")
			continue
		}

		fileInfos = append(fileInfos, models.FileInfo{
			Name:          file.Name(),
			Path:          filepath.ToSlash(filePath),
			IsDir:         file.IsDir(),
			FileSize:      helper.FormatFileSize(info, file),
			Size:          info.Size(),
			LastModified:  info.ModTime().Format(timeFormat),
			IsEditable:    helper.IsText(filepath.ToSlash(filePath)),
			Permissions:   fm.getFilePermissions(info),
			FileType:      fm.getFileType(info),
			Owner:         owner,
			Group:         group,
			CreationDate:  fm.getCreationDate(info),
		})
	}

	helper.SortFileInfos(fileInfos)
	return fileInfos, nil
}

// getFilePermissions returns the file permissions in octal format
func (fm *FileManager) getFilePermissions(fi os.FileInfo) string {
	return fmt.Sprintf("%04o", fi.Mode().Perm())
}

// getFileType returns the type of file
func (fm *FileManager) getFileType(fi os.FileInfo) string {
	if fi.IsDir() {
		return "directory"
	}
	return "file"
}

// getFileOwnerGroup returns the file owner and group for Unix systems
func (fm *FileManager) getFileOwnerGroup(path string) (string, string, error) {
	cmd := exec.Command("stat", "-c", "%U:%G", path)
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get file owner and group: %w", err)
	}

	parts := strings.SplitN(strings.TrimSpace(string(output)), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected output format from stat command")
	}

	return parts[0], parts[1], nil
}

// getCreationDate returns the file creation date (or modification date if creation date is not available)
func (fm *FileManager) getCreationDate(fi os.FileInfo) string {
	return fi.ModTime().Format(timeFormat)
}

// ListDirectory lists the contents of a directory
func (fm *FileManager) ListDirectory(path string) ([]models.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return fm.PrepareFileInfo(entries, path)
}

// CreateDirectory creates a new directory
func (fm *FileManager) CreateDirectory(path string) error {
	if !strings.HasPrefix(path, fm.Config.Files.StorageDir) {
		return fmt.Errorf("cannot create directory outside of storage area")
	}
	if err := os.MkdirAll(path, defaultFilePermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// DeleteFile deletes a file or empty directory
func (fm *FileManager) DeleteFile(path string) error {
	if !strings.HasPrefix(path, fm.Config.Files.StorageDir) {
		return fmt.Errorf("cannot delete file outside of storage area")
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// RenameFile renames a file or directory
func (fm *FileManager) RenameFile(oldPath, newPath string) error {
	if !strings.HasPrefix(oldPath, fm.Config.Files.StorageDir) || !strings.HasPrefix(newPath, fm.Config.Files.StorageDir) {
		return fmt.Errorf("cannot rename file outside of storage area")
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return nil
}

// ValidateParams checks the validity of the creation type and name.
func ValidateParams(creationType, name string) error {
	if creationType == "" || (creationType != "file" && creationType != "dir") {
		return fmt.Errorf("invalid 'type' parameter. must be 'file' or 'dir'")
	}
	if name == "" {
		return fmt.Errorf("name parameter must be provided")
	}
	return nil
}

// ResolvePath returns the cleaned and defaulted current path.
func ResolvePath(path string) string {
	if path == "" {
		return "." // Default to the current working directory
	}
	return filepath.Clean(path)
}

// IsValidPath checks if the new path is within the base directory.
func IsValidPath(baseDir, newPath string) bool {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	absNewPath, err := filepath.Abs(newPath)
	if err != nil {
		log.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	return strings.HasPrefix(absNewPath, absBaseDir)
}

// CreateEntity creates a file or directory based on the type and path provided.
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