package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RenameHandler menangani permintaan untuk mengganti nama file
func RenameHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		// return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
		return respondWithError(c,fiber.StatusMethodNotAllowed,"Method not allowed")
	}

	var payload struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		return respondWithError(c, fiber.StatusBadRequest, "Failed to decode request payload: "+ err.Error())
	}

	// Membersihkan dan memvalidasi jalur lama
	oldFilePath := filepath.Clean(payload.OldPath)
	if !utils.IsValidPath(oldFilePath) || strings.Contains(payload.NewName, "..") {
		log.Printf("Invalid path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid path")
	}

	// Mendapatkan ekstensi dari file lama
	extension := filepath.Ext(oldFilePath)

	// Menambahkan ekstensi jika tidak ada
	if filepath.Ext(payload.NewName) == "" {
		payload.NewName += extension
	}

	// Membentuk jalur baru
	newPath := filepath.Join(filepath.Dir(oldFilePath), payload.NewName)

	// Memeriksa apakah file lama ada
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		log.Printf("File does not exist: %v", err)
		// return c.Status(fiber.StatusNotFound).SendString("File does not exist")
		return respondWithError(c, fiber.StatusBadRequest, "File does not exist: "+ err.Error())
	}

	// Mencoba mengganti nama file
	if err := os.Rename(oldFilePath, newPath); err != nil {
		log.Printf("Failed to rename file: %v", err)
		// return c.Status(fiber.StatusInternalServerError).SendString("Failed to rename file")
		return respondWithError(c, fiber.StatusBadRequest, "Failed to "+err.Error())
	}

	log.Println("File renamed successfully")
	// return c.SendStatus(fiber.StatusNoContent)
	return respondWithJSON(c,fiber.StatusOK,models.Response{
		Message:"File renamed successfully",
	})
}
