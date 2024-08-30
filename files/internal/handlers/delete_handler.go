package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// DeleteHandler menangani permintaan untuk menghapus file atau direktori
func DeleteHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodDelete {
		// return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		return respondWithError(c,fiber.StatusMethodNotAllowed,"Method not allowed")
	}

	var payload struct {
		Path string `json:"path"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		return respondWithError(c, fiber.StatusBadRequest, "Failed to decode request payload: "+ err.Error())
	}

	if !utils.IsValidPath(payload.Path) {
		log.Printf("Invalid path")
		// return respondWithError(c, fiber.StatusBadRequest, "Invalid path")
		return respondWithError(c,fiber.StatusBadRequest,"Invalid path")
	}

	if err := os.RemoveAll(payload.Path); err != nil {
		log.Printf("Failed to delete file: %v", err)
		// return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete file")
		return respondWithError(c,fiber.StatusInternalServerError,"Failed to delete file: "+err.Error())
	}

	log.Println("File deleted successfully")
	// return c.SendStatus(fiber.StatusNoContent)
	return respondWithJSON(c,fiber.StatusOK,models.Response{
		Message:"File deleted successfully",
	})
}
