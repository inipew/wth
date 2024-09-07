package handlers

import (
	"files/internal/models"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func UploadFileHandler(c *fiber.Ctx) error {
	// Pastikan metode adalah POST
	if c.Method() != fiber.MethodPost {
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Invalid request method")
	}

	// Ambil file yang diunggah
	file, err := c.FormFile("file")
	if err != nil {
		return respondWithError(c, fiber.StatusBadRequest, "Unable to retrieve file: " + err.Error())

    }

	// Tentukan path tujuan
	destPath := c.FormValue("path")
	log.Printf("path: %s", destPath)
	if destPath == "" {
		destPath = "./uploads" // Default directory
	}

	// Sanitasi path dan filename
	destPath = filepath.Clean(destPath)
	filePath := filepath.Join(destPath, file.Filename)
	if strings.Contains(file.Filename, "..") {
		return respondWithError(c, fiber.StatusBadRequest, "Invalid filename")
	}

	// Pastikan direktori tujuan ada
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		log.Logger.Error().Err(err).Msg("Unable to create directory")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to create directory: "+err.Error())
	}

	// Buka file yang diunggah
	src, err := file.Open()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Unable to open file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to open file: "+err.Error())

	}
	defer src.Close()

	// Buat file baru
	dst, err := os.Create(filePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Unable to create file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to create file: "+ err.Error())

	}
	defer dst.Close()

	// Salin isi file
	if _, err := io.Copy(dst, src); err != nil {
		log.Logger.Error().Err(err).Msg("Unable to create file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to save file: "+err.Error())
	}

	// return c.SendString("File uploaded successfully")
	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message:"File uploaded successfully",
	})
}
