package handlers

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func DownloadHandler(c *fiber.Ctx) error {
	fileParam := c.Query("file")
	if fileParam == "" {
		// return c.Status(fiber.StatusBadRequest).SendString("File parameter is required")
		return respondWithJSON(c, fiber.StatusBadRequest,"File parameter is required")

	}

	// Clean and get the absolute path of the file
	absFilePath, err := filepath.Abs(filepath.Clean(fileParam))
	if err != nil {
		// return c.Status(fiber.StatusBadRequest).SendString("Invalid file path")
		return respondWithJSON(c, fiber.StatusBadRequest,"Invalid file path"+err.Error())
	}

	// Set headers
	c.Set("Content-Disposition", "attachment; filename="+filepath.Base(absFilePath))
	c.Set("Content-Type", "application/octet-stream")

	// Serve the file
	return c.SendFile(absFilePath)
}
