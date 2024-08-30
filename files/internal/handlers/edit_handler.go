package handlers

import (
	"bufio"
	"bytes"
	"files/internal/models"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// FileContentRequest digunakan untuk format permintaan untuk konten file
type FileContentRequest struct {
	FileName string `json:"fileName"`
	Content  string `json:"content"`
}

// ViewHandler menangani permintaan untuk melihat konten file
func ViewHandler(c *fiber.Ctx) error {
    fileName := c.Query("filepath")
    if fileName == "" {
        // return c.Status(fiber.StatusBadRequest).SendString("File path is required")
		return respondWithJSON(c, fiber.StatusBadRequest,"File path is required")

    }

    decodedFileName, err := url.QueryUnescape(fileName)
    if err != nil {
        // return c.Status(fiber.StatusBadRequest).SendString("Invalid file path")
		return respondWithJSON(c, fiber.StatusBadRequest,"Invalid filename path: "+err.Error())
    }
    
    filePath, err := filepath.Abs(filepath.Clean(decodedFileName))
    if err != nil {
        // return c.Status(fiber.StatusBadRequest).SendString("Invalid file path")
		return respondWithJSON(c, fiber.StatusBadRequest,"Invalid file path: "+err.Error())
    }

    file, err := os.Open(filePath)
    if err != nil {
        // return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
		return respondWithJSON(c, fiber.StatusInternalServerError,"Failed to open file: "+err.Error())
    }
    defer file.Close()

    reader := bufio.NewReader(file)
	var content bytes.Buffer
	if _, err := io.Copy(&content, reader); err != nil {
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
		return respondWithJSON(c, fiber.StatusInternalServerError,"Failed to read file: "+err.Error())
	}

    response := FileContentRequest{
        FileName: filePath,
        Content:  content.String(),
    }

    // return c.Status(fiber.StatusOK).JSON(response)
	return respondWithJSON(c,fiber.StatusOK,response)
}


// SaveHandler menangani permintaan untuk menyimpan konten file
func SaveHandler(c *fiber.Ctx) error {
	var req FileContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(req.FileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file path")
	}

	// Write the content back to the file
	err = os.WriteFile(absPath, []byte(req.Content), 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	// return c.SendStatus(fiber.StatusOK)
	return respondWithJSON(c, fiber.StatusOK,models.Response{Message: "File Saved"})
}
