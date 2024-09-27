package handlers

import (
	"files/internal/core/archive"
	"files/internal/models"
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ArchiveHandler handles view archive file
func (h *Handlers) ArchiveHandler(c *fiber.Ctx) error {
	return h.handleFileOperation(c, h.viewHandler)
}

// ArchiveHandler handles view archive file
func (h *Handlers) ExtractorHandler(c *fiber.Ctx) error {
	return h.handleFileOperation(c, h.unzipHandler)
}

func (h *Handlers) viewHandler(c *fiber.Ctx) error {
	archivePath := c.Query("path")
	if archivePath == "" {
		return models.RespondWithError(c, fiber.StatusBadRequest, "Missing 'path' query parameter")
	}

	archivePath = filepath.Clean(archivePath)
	fileInfos, err := archive.ProcessArchiveFile(archivePath)
	if err != nil {
		log.Error().Err(err).Str("path", archivePath).Msg("Failed to process archive file")
		return models.RespondWithError(c, fiber.StatusInternalServerError, fmt.Sprintf("Error processing archive file: %v", err))
	}

	return c.Status(fiber.StatusOK).JSON(models.ArchiveInfo{
		Name:  filepath.Base(archivePath),
		Path:  archivePath,
		Files: fileInfos,
	})
}

// UnzipHandler handles the extraction of various archive formats
func (h *Handlers) unzipHandler(c *fiber.Ctx) error {
	filePath, err := archive.GetAndValidateFilePath(c.Query("file"))
	if err != nil {
		return models.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	extractor, err := archive.GetExtractor(filePath)
	if err != nil {
		return models.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := extractor(filePath); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("Failed to extract file")
		return models.RespondWithError(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to extract file: %v", err))
	}

	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File extracted successfully",
	})
}