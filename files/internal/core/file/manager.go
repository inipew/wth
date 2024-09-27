package file

import (
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

// getDirectoryPath extracts and validates the directory path from the request
func GetDirectoryPath(c *fiber.Ctx) (string, error) {
	dir := c.Query("path")
	if dir == "" {
		return os.Getwd()
	}

	decodedDir, err := url.QueryUnescape(dir)
	if err != nil {
		return "", err
	}
	cleanedDir := filepath.Clean(decodedDir)
	return filepath.Abs(cleanedDir)
}

// prepareFileInfo prepares FileInfo slice from directory entries
func PrepareFileInfo(files []os.DirEntry, dirPath string) ([]models.FileInfo, error) {
	var fileInfos []models.FileInfo

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Logger.Error().Err(err).Str("file", file.Name()).Msg("Error getting file info")
			continue
		}
		owner, group, _ := getFileOwnerGroup(filepath.Join(dirPath, file.Name()))
		fileInfos = append(fileInfos, models.FileInfo{
			Name:          file.Name(),
			Path:          filepath.ToSlash(filepath.Join(dirPath, file.Name())),
			IsDir:         file.IsDir(),
			FileSize:      helper.FormatFileSize(info, file),
			Size: 		   info.Size(),
			LastModified:  info.ModTime().Format("2006-01-02 15:04:05"),
			IsEditable:    helper.IsText(filepath.ToSlash(filepath.Join(dirPath, file.Name()))),
			Permissions:   getFilePermissions(info),
			FileType:      getFileType(info),
			Owner:         owner,
			Group:         group,
			CreationDate:  getCreationDate(info),
		})
	}
	helper.SortFileInfos(fileInfos)
	return fileInfos, nil
}

// getFilePermissions returns the file permissions in octal format
func getFilePermissions(fi os.FileInfo) string {
	return fmt.Sprintf("%04o", fi.Mode().Perm())
}

// getFileType returns the type of file
func getFileType(fi os.FileInfo) string {
	if fi.IsDir() {
		return "directory"
	}
	return "file"
}

// getFileOwnerGroup returns the file owner and group for Unix systems
func getFileOwnerGroup(path string) (string, string, error) {
	cmd := exec.Command("stat", "-c", "%U:%G", path)
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	parts := strings.SplitN(strings.TrimSpace(string(output)), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected output format from stat command")
	}

	return parts[0], parts[1], nil
}

// getCreationDate returns the file creation date (or modification date if creation date is not available)
func getCreationDate(fi os.FileInfo) string {
	return fi.ModTime().Format("2006-01-02 15:04:05")
}
