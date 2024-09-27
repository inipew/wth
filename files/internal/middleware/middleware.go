package middleware

import (
	"files/internal/utils/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS configures and applies CORS middleware to the Fiber app
func SetupCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
}

// SetupCompression configures and applies compression middleware to the Fiber app
func SetupCompression(app *fiber.App) {
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}

// RequestLogger returns a middleware function that logs request details
func RequestLogger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		log.Info(
			"Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", duration,
		)

		return err
	}
}