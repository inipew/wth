package handlers

import (
	"files/internal/models"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFileHandler(c *fiber.Ctx) error {
	// Pastikan metode adalah POST
	if c.Method() != fiber.MethodPost {
		// return fiber.NewError(fiber.StatusMethodNotAllowed, "Invalid request method")
		return respondWithError(c, fiber.StatusMethodNotAllowed, "Invalid request method")
	}

	// Ambil file yang diunggah
	file, err := c.FormFile("file")
	if err != nil {
        // return fiber.NewError(fiber.StatusBadRequest, "Unable to retrieve file: " + err.Error())
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
		// return fiber.NewError(fiber.StatusBadRequest, "Invalid filename")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid filename")

	}

	// Pastikan direktori tujuan ada
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		// return fiber.NewError(fiber.StatusInternalServerError, "Unable to create directory")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to create directory: "+err.Error())

	}

	// Buka file yang diunggah
	src, err := file.Open()
	if err != nil {
		// return fiber.NewError(fiber.StatusInternalServerError, "Unable to open file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to open file: "+err.Error())

	}
	defer src.Close()

	// Buat file baru
	dst, err := os.Create(filePath)
	if err != nil {
		// return fiber.NewError(fiber.StatusInternalServerError, "Unable to create file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to create file: "+ err.Error())

	}
	defer dst.Close()

	// Salin isi file
	if _, err := io.Copy(dst, src); err != nil {
		// return fiber.NewError(fiber.StatusInternalServerError, "Unable to save file")
		return respondWithError(c, fiber.StatusInternalServerError, "Unable to save file: "+err.Error())
	}

	// return c.SendString("File uploaded successfully")
	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message:"File uploaded successfully",
	})
}
