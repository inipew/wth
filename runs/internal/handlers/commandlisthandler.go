package handlers

import (
	"runs/internal/config"

	"github.com/gofiber/fiber/v2"
)

// getCommandList handles the API request for listing all commands
func GetCommandList(c *fiber.Ctx) error {
	config.ConfigMutex.RLock()
	defer config.ConfigMutex.RUnlock()

	if config.ConfigData.Commands == nil {
		return respondWithError(c, fiber.StatusNotFound, "No commands available")
	}

	return respondWithJSON(c, fiber.StatusOK, config.ConfigData.Commands)
}