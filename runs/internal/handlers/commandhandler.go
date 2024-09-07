package handlers

import (
	"context"
	"fmt"
	cfg "runs/internal/config"
	"runs/internal/models"
	"runs/internal/utils"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const defaultTimeout = 10 * time.Second

// CommandHandler handles the command execution
func CommandHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), defaultTimeout)
	defer cancel()

	value := c.FormValue("value")
	if value == "" {
		return handleError(c, fiber.StatusBadRequest, "Missing 'value' parameter")
	}

	command, err := processCommand(c, value)
	if err != nil {
		return handleError(c, fiber.StatusBadRequest, err.Error())
	}

	output, err := utils.RunCommand(ctx, command.Command, defaultTimeout)
	if err != nil {
		log.Logger.Error().Str("command", command.Command).Err(err).Msg("Failed to execute command")
		return handleError(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to execute command '%s': %v", command.Command, err))
	}

	return respondWithJSON(c, fiber.StatusOK, models.Response{Message: output})
}

// processCommand processes the command based on the given value
func processCommand(c *fiber.Ctx, value string) (*cfg.Command, error) {
	if value == "custom" {
		customCommand := c.FormValue("custom_command")
		if customCommand == "" {
			return nil, fmt.Errorf("missing 'custom_command' parameter")
		}
		return &cfg.Command{Command: customCommand}, nil
	}
	command, err := cfg.FindCommandByValue(value)
	if err != nil {
		return nil, fmt.Errorf("command with Value '%s' not found", value)
	}
	command.Command = cfg.ReplacePlaceholders(command.Command, nil)
	return command, nil
}

// handleError logs and responds with an error message
func handleError(c *fiber.Ctx, status int, message string) error {
	log.Logger.Error().Str("status", fmt.Sprintf("%d", status)).Msg(message)
	return respondWithError(c, status, message)
}

// respondWithJSON marshals payload to JSON and sends it as a response
func respondWithJSON(c *fiber.Ctx, status int, payload any) error {
	data, err := sonic.Marshal(payload)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating JSON response")
		return respondWithError(c, fiber.StatusInternalServerError, "Error generating JSON response")
	}
	return c.Status(status).Send(data)
}

// respondWithError responds with an error message as JSON
func respondWithError(c *fiber.Ctx, status int, message string) error {
	response := models.Response{Message: message}
	data, err := sonic.Marshal(response)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating error response")
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating JSON response")
	}
	return c.Status(status).Send(data)
}
