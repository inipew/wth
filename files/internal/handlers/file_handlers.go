package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// FileHandler handles file listing requests
func FileHandler(c *fiber.Ctx) error {
	currentPath, err := getDirectoryPath(c)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get directory path")
		return respondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	if !utils.IsValidPath(currentPath) {
		errMsg := "Invalid path"
		log.Warn().Str("path", currentPath).Msg(errMsg)
		return respondWithError(c, fiber.StatusBadRequest, errMsg)
	}

	files, err := os.ReadDir(currentPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", currentPath).Msg("Error reading directory")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to read directory: "+err.Error())
	}

	fileInfos, err := prepareFileInfo(files, currentPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", currentPath).Msg("Failed to get file info")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to get file info: "+err.Error())
	}

	previousPath := filepath.Dir(currentPath)
	if currentPath == previousPath {
		previousPath = ""
	}

	return respondWithJSON(c, fiber.StatusOK, models.DirectoryInfo{
		CurrentPath:  currentPath,
		PreviousPath: previousPath,
		Files:        fileInfos,
	})
}

// getDirectoryPath extracts and validates the directory path from the request
func getDirectoryPath(c *fiber.Ctx) (string, error) {
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
func prepareFileInfo(files []os.DirEntry, dirPath string) ([]models.FileInfo, error) {
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
			FileSize:      utils.FormatFileSize(info, file),
			Size: 		   info.Size(),
			LastModified:  info.ModTime().Format("2006-01-02 15:04:05"),
			IsEditable:    utils.IsText(filepath.ToSlash(filepath.Join(dirPath, file.Name()))),
			Permissions:   getFilePermissions(info),
			FileType:      getFileType(info),
			Owner:         owner,
			Group:         group,
			CreationDate:  getCreationDate(info),
		})
	}
	utils.SortFileInfos(fileInfos)
	return fileInfos, nil
}

// respondWithJSON creates a JSON response using Sonic
func respondWithJSON(c *fiber.Ctx, status int, payload any) error {
	data, err := sonic.Marshal(payload)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating JSON response")
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating JSON response")
	}
	return c.Status(status).Send(data)
}

// respondWithError creates an error response using Sonic
func respondWithError(c *fiber.Ctx, status int, message string) error {
	response := models.Response{
		Message: message,
	}
	data, err := sonic.Marshal(response)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating JSON response")
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating JSON response")
	}
	return c.Status(status).Send(data)
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
