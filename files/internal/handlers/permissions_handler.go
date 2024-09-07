package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// UpdatePermissionsHandler handles updating file permissions
func UpdatePermissionsHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPut {
		log.Logger.Warn().Str("method", c.Method()).Msg("Method not allowed")
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}

	var payload struct {
		Path        string `json:"path"`
		Permissions string `json:"permissions"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to decode request payload")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	if !utils.IsValidPath(payload.Path) {
		log.Logger.Warn().Str("path", payload.Path).Msg("Invalid file path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid file path")
	}

	if _, err := os.Stat(payload.Path); os.IsNotExist(err) {
		log.Logger.Error().Err(err).Str("path", payload.Path).Msg("File not found")
		return respondWithError(c, fiber.StatusNotFound, "File not found: "+err.Error())
	}

	permissions, err := strconv.ParseInt(payload.Permissions, 8, 32)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse permissions")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid permissions format: "+err.Error())
	}

	fileMode := os.FileMode(permissions)

	if err := os.Chmod(payload.Path, fileMode); err != nil {
		log.Logger.Error().Err(err).Str("path", payload.Path).Msg("Failed to update file permissions")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to update file permissions: "+err.Error())
	}

	log.Logger.Info().Str("path", payload.Path).Msg("File permissions updated successfully")
	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File permissions updated successfully",
	})
}
