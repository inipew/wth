package handlers

import (
	"files/internal/models"
	"files/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RenameHandler menangani permintaan untuk mengganti nama file
func RenameHandler(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		log.Logger.Warn().Str("method", c.Method()).Msg("Method not allowed")
		return respondWithError(c,fiber.StatusMethodNotAllowed,"Method not allowed")
	}

	var payload struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to decode request payload")
		return respondWithError(c, fiber.StatusBadRequest, "Failed to decode request payload: "+ err.Error())
	}

	// Membersihkan dan memvalidasi jalur lama
	oldFilePath := filepath.Clean(payload.OldPath)
	if !utils.IsValidPath(oldFilePath) || strings.Contains(payload.NewName, "..") {
		log.Logger.Warn().Str("path", oldFilePath).Msg("Invalid file path")
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
		log.Logger.Error().Err(err).Str("path", oldFilePath).Msg("File not found")
		return respondWithError(c, fiber.StatusBadRequest, "File does not exist: "+ err.Error())
	}

	// Mencoba mengganti nama file
	if err := os.Rename(oldFilePath, newPath); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to rename file")
		return respondWithError(c, fiber.StatusBadRequest, "Failed to "+err.Error())
	}
	// return c.SendStatus(fiber.StatusNoContent)
	return respondWithJSON(c,fiber.StatusOK,models.Response{
		Message:"File renamed successfully",
	})
}
