package models

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Response digunakan untuk format balasan API
type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// respondWithJSON creates a JSON response using Sonic
func RespondWithJSON(c *fiber.Ctx, status int, payload any) error {
	data, err := sonic.Marshal(payload)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating JSON response")
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating JSON response")
	}
	return c.Status(status).Send(data)
}

// respondWithError creates an error response using Sonic
func RespondWithError(c *fiber.Ctx, status int, message string) error {
	response := Response{
		Message: message,
	}
	data, err := sonic.Marshal(response)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error generating JSON response")
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating JSON response")
	}
	return c.Status(status).Send(data)
}