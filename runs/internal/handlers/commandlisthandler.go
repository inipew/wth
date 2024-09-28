package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// GetCommandList handles the API request for listing all commands
func (h *Handler) GetCommandList(c *fiber.Ctx) error {
	cfg := h.configManager.GetConfig()
	if cfg == nil {
		return respondWithError(c, fiber.StatusInternalServerError, "Configuration not loaded")
	}

	if len(cfg.Commands) == 0 {
		return respondWithError(c, fiber.StatusNotFound, "No commands available")
	}

	return respondWithJSON(c, fiber.StatusOK, cfg.Commands)
}