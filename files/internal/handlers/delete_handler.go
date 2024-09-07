package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// DeleteHandler handles requests to delete files or directories
func DeleteHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodDelete {
		log.Logger.Warn().Msg("Method not allowed for delete operation")
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}

	var payload struct {
		Path string `json:"path"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to decode request payload")
		return respondWithError(c, fiber.StatusBadRequest, "Failed to decode request payload: "+err.Error())
	}

	if !utils.IsValidPath(payload.Path) {
		log.Logger.Warn().Str("path", payload.Path).Msg("Invalid path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid path")
	}

	if err := os.RemoveAll(payload.Path); err != nil {
		log.Logger.Error().Err(err).Str("path", payload.Path).Msg("Failed to delete file")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to delete file: "+err.Error())
	}

	log.Logger.Info().Str("path", payload.Path).Msg("File deleted successfully")
	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File deleted successfully",
	})
}
