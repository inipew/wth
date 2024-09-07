package handlers

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func DownloadHandler(c *fiber.Ctx) error {
	fileParam := c.Query("file")
	if fileParam == "" {
		log.Logger.Warn().Msg("File parameter is missing")
		return respondWithJSON(c, fiber.StatusBadRequest, "File parameter is required")
	}

	// Clean and get the absolute path of the file
	absFilePath, err := filepath.Abs(filepath.Clean(fileParam))
	if err != nil {
		log.Logger.Error().Err(err).Str("fileParam", fileParam).Msg("Failed to resolve file path")
		return respondWithJSON(c, fiber.StatusBadRequest, "Invalid file path: "+err.Error())
	}

	// Check if the file exists
	if _, err := filepath.Abs(absFilePath); err != nil {
		log.Logger.Error().Err(err).Str("absFilePath", absFilePath).Msg("Failed to access file")
		return respondWithJSON(c, fiber.StatusNotFound, "File not found: "+err.Error())
	}

	// Set headers
	c.Set("Content-Disposition", "attachment; filename="+filepath.Base(absFilePath))
	c.Set("Content-Type", "application/octet-stream")

	// Serve the file
	err = c.SendFile(absFilePath)
	if err != nil {
		log.Logger.Error().Err(err).Str("filePath", absFilePath).Msg("Failed to send file")
		return respondWithJSON(c, fiber.StatusInternalServerError, "Failed to send file: "+err.Error())
	}

	log.Logger.Info().Str("filePath", absFilePath).Msg("File served successfully")
	return nil
}
