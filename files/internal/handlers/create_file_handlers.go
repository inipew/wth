package handlers

import (
	"files/internal/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// MakeNewHandler handles requests to create a new file or directory based on query parameters.
func MakeNewHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		log.Logger.Warn().Msg("Invalid request method for MakeNewHandler")
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Invalid request method")
	}

	// Extract parameters from query
	creationType := c.Query("type")
	currentPath := c.Query("currentPath")
	name := c.Query("name")

	// Validate parameters
	if err := validateParams(creationType, name); err != nil {
		log.Logger.Error().Err(err).Msg("Parameter validation failed")
		return respondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	// Resolve and validate the current path
	currentPath = resolvePath(currentPath)
	baseDir, err := os.Getwd()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get the base directory")
		return respondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	if !isValidPath(baseDir, currentPath) {
		log.Logger.Warn().Str("currentPath", currentPath).Msg("Invalid current path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid current path")
	}

	// Create file or directory
	if err := createEntity(creationType, currentPath, name); err != nil {
		log.Logger.Error().Err(err).Str("type", creationType).Str("path", currentPath).Msg("Error creating entity")
		return respondWithError(c, fiber.StatusInternalServerError, err.Error())
	}

	log.Logger.Info().Str("type", creationType).Str("path", filepath.Join(currentPath, name)).Msg("File or directory created successfully")
	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message: fmt.Sprintf("%s created successfully", creationType),
	})
}

// validateParams checks the validity of the creation type and name.
func validateParams(creationType, name string) error {
	if creationType == "" || (creationType != "file" && creationType != "dir") {
		return fmt.Errorf("invalid 'type' parameter. must be 'file' or 'dir'")
	}
	if name == "" {
		return fmt.Errorf("name parameter must be provided")
	}
	return nil
}

// resolvePath returns the cleaned and defaulted current path.
func resolvePath(path string) string {
	if path == "" {
		return "." // Default to the current working directory
	}
	return filepath.Clean(path)
}

// isValidPath checks if the new path is within the base directory.
func isValidPath(baseDir, newPath string) bool {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	absNewPath, err := filepath.Abs(newPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("baseDir", baseDir).Str("newPath", newPath).Msg("Failed to get absolute path")
		return false
	}
	return strings.HasPrefix(absNewPath, absBaseDir)
}

// createEntity creates a file or directory based on the type and path provided.
func createEntity(creationType, currentPath, name string) error {
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
