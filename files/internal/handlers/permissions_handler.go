package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdatePermissionsHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPut {
		// return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Method not allowed")
	}

	var payload struct {
		Path        string `json:"path"`
		Permissions string `json:"permissions"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		// return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid request payload: "+err.Error())

	}

	if !utils.IsValidPath(payload.Path) {
		// return c.Status(fiber.StatusBadRequest).SendString("Invalid file path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid file path")

	}

	// Periksa apakah file ada
	if _, err := os.Stat(payload.Path); os.IsNotExist(err) {
		// return c.Status(fiber.StatusNotFound).SendString("File not found")
		return respondWithError(c, fiber.StatusNotFound, "File not found: "+err.Error())

	}

	// Parsing izin
	permissions, err := strconv.ParseInt(payload.Permissions, 8, 32)
	if err != nil {
		log.Printf("Failed to parse permissions: %v", err)
		// return c.Status(fiber.StatusBadRequest).SendString("Invalid permissions format")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid permissions format: "+err.Error())

	}

	fileMode := os.FileMode(permissions)

	// Menerapkan izin
	if err := os.Chmod(payload.Path, fileMode); err != nil {
		log.Printf("Failed to update file permissions: %v", err)
		// return c.Status(fiber.StatusInternalServerError).SendString("Failed to update file permissions")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to update file permissions: " + err.Error())
	}

	log.Println("File permissions updated successfully")
	// return c.SendStatus(fiber.StatusNoContent)
	return respondWithJSON(c,fiber.StatusOK,models.Response{
		Message:"File permissions updated successfully",
	})
}
