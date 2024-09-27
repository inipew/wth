package handlers

import (
	"bufio"
	"bytes"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"files/internal/core/file"
	"files/internal/models"
	"files/internal/utils/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// FileHandler handles file listing requests
func FileHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, listFiles)
}

// DeleteHandler handles file deletion requests
func DeleteHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, deleteFile)
}

// DownloadHandler handles file download requests
func DownloadHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, downloadFile)
}

// RenameHandler handles file renaming requests
func RenameHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, renameFile)
}

// UploadFileHandler handles file upload requests
func UploadFileHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, uploadFile)
}

// ViewHandler handles requests to view file content
func ViewHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, viewFile)
}

// SaveHandler handles requests to save file content
func SaveHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, saveFile)
}

// MakeNewHandler handles requests to create new files or directories
func MakeNewHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, makeNew)
}

// FileOperation represents a function that performs a file operation
type FileOperation func(*fiber.Ctx) error

// handleFileOperation is a generic handler for file operations
func handleFileOperation(c *fiber.Ctx, operation FileOperation) error {
	if err := operation(c); err != nil {
		status, message := parseError(err)
		return respondWithError(c, status, message)
	}
	return nil
}

// listFiles handles listing files in a directory
func listFiles(c *fiber.Ctx) error {
	currentPath, err := file.GetDirectoryPath(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to get directory path: "+err.Error())
	}

	if !helper.IsValidPath(currentPath) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid path")
	}

	files, err := os.ReadDir(currentPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read directory: "+err.Error())
	}

	fileInfos, err := file.PrepareFileInfo(files, currentPath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get file info: "+err.Error())
	}

	previousPath := filepath.Dir(currentPath)
	if currentPath == previousPath {
		previousPath = ""
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.DirectoryInfo{
		CurrentPath:  currentPath,
		PreviousPath: previousPath,
		Files:        fileInfos,
	})
}

// deleteFile handles file deletion
func deleteFile(c *fiber.Ctx) error {
	var payload struct {
		Path string `json:"path"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to decode request payload: "+err.Error())
	}

	if !helper.IsValidPath(payload.Path) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid path")
	}

	if err := os.RemoveAll(payload.Path); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete file: "+err.Error())
	}

	log.Info().Str("path", payload.Path).Msg("File deleted successfully")
	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File deleted successfully",
	})
}

// downloadFile handles file download
func downloadFile(c *fiber.Ctx) error {
	fileParam := c.Query("file")
	if fileParam == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File parameter is required")
	}

	absFilePath, err := filepath.Abs(filepath.Clean(fileParam))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file path: "+err.Error())
	}

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found: "+err.Error())
	}

	return c.Download(absFilePath)
}

// renameFile handles file renaming
func renameFile(c *fiber.Ctx) error {
	var payload struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to decode request payload: "+err.Error())
	}

	oldFilePath := filepath.Clean(payload.OldPath)
	if !helper.IsValidPath(oldFilePath) || strings.Contains(payload.NewName, "..") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid path")
	}

	newPath := filepath.Join(filepath.Dir(oldFilePath), payload.NewName)

	if err := os.Rename(oldFilePath, newPath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to rename file: "+err.Error())
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File renamed successfully",
	})
}

// uploadFile handles file upload
func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to retrieve file: "+err.Error())
	}

	destPath := c.FormValue("path")
	destPath = filepath.Clean(destPath)

	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to create directory: "+err.Error())
	}

	filePath := filepath.Join(destPath, file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to save file: "+err.Error())
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File uploaded successfully",
	})
}

// FileContentRequest represents the request format for file content
type FileContentRequest struct {
	FileName string `json:"fileName"`
	Content  string `json:"content"`
}

// viewFile handles viewing file content
func viewFile(c *fiber.Ctx) error {
	fileName := c.Query("filepath")
	if fileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "File path is required")
	}

	decodedFileName, err := url.QueryUnescape(fileName)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid filename path: "+err.Error())
	}

	filePath, err := filepath.Abs(filepath.Clean(decodedFileName))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file path: "+err.Error())
	}

	content, err := readFileContent(filePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read file: "+err.Error())
	}

	return models.RespondWithJSON(c, fiber.StatusOK, FileContentRequest{
		FileName: filePath,
		Content:  content,
	})
}

// saveFile handles saving file content
func saveFile(c *fiber.Ctx) error {
	var req FileContentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	absPath, err := filepath.Abs(req.FileName)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file path: "+err.Error())
	}

	if err := os.WriteFile(absPath, []byte(req.Content), 0644); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file: "+err.Error())
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{Message: "File Saved"})
}

// makeNew handles creating new files or directories
func makeNew(c *fiber.Ctx) error {
	creationType := c.Query("type")
	currentPath := c.Query("currentPath")
	name := c.Query("name")

	if err := file.ValidateParams(creationType, name); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	currentPath = file.ResolvePath(currentPath)
	baseDir, err := os.Getwd()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get the base directory: "+err.Error())
	}

	if !file.IsValidPath(baseDir, currentPath) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid current path")
	}

	if err := file.CreateEntity(creationType, currentPath, name); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating entity: "+err.Error())
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: creationType + " created successfully",
	})
}

// readFileContent is a helper function to read file content
func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content bytes.Buffer
	if _, err := io.Copy(&content, reader); err != nil {
		return "", err
	}

	return content.String(), nil
}

// parseError is a helper function to parse errors and return appropriate status codes
func parseError(err error) (int, string) {
	if e, ok := err.(*fiber.Error); ok {
		return e.Code, e.Message
	}
	return fiber.StatusInternalServerError, err.Error()
}

// respondWithError is a helper function to respond with an error
func respondWithError(c *fiber.Ctx, status int, message string) error {
	log.Error().Int("status", status).Msg(message)
	return models.RespondWithError(c, status, message)
}