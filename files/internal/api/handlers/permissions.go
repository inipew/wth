package handlers

import (
	"files/internal/models"
	"files/internal/utils/helper"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// UpdatePermissionsHandler handles updating file permissions
func UpdatePermissionsHandler(c *fiber.Ctx) error {
	return handleFileOperation(c, updatePermissions)
}

// updatePermissions is the core function to update file permissions
func updatePermissions(c *fiber.Ctx) error {
	var payload struct {
		Path        string `json:"path"`
		Permissions string `json:"permissions"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	if err := validateUpdatePermissionsInput(payload.Path, payload.Permissions); err != nil {
		return err
	}

	fileMode, err := parsePermissions(payload.Permissions)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid permissions format: "+err.Error())
	}

	if err := os.Chmod(payload.Path, fileMode); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update file permissions: "+err.Error())
	}

	log.Info().Str("path", payload.Path).Msg("File permissions updated successfully")
	return models.RespondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File permissions updated successfully",
	})
}

// validateUpdatePermissionsInput validates the input for updating permissions
func validateUpdatePermissionsInput(path, permissions string) error {
	if !helper.IsValidPath(path) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found: "+err.Error())
	}

	if _, err := strconv.ParseInt(permissions, 8, 32); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid permissions format: "+err.Error())
	}

	return nil
}

// parsePermissions converts a string representation of permissions to os.FileMode
func parsePermissions(permissions string) (os.FileMode, error) {
	p, err := strconv.ParseInt(permissions, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(p), nil
}