package main

import (
	"embed"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"files/internal/handlers"
	"files/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	if err := logger.InitLogger(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	// Load configuration
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to load config file")
	}

	port := cfg.Section("server").Key("port").String()
	if port == "" {
		log.Logger.Fatal().Msg("Port not defined in config.ini")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:                  errorHandler,
		DisablePreParseMultipartForm:  true,
		StreamRequestBody:             true,
	})

	// Middleware
	setupMiddleware(app)
	// API routes
	defineAPIRoutes(app)

	// // Serve static files
	// app.Static("/", "./frontend/dist", fiber.Static{
	// 	Compress: true,
	//  CacheDuration: 3 * time.Hour,
	// })


	// // Fallback route for SPA
	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	if _, err := os.Stat(filepath.Join("./frontend/dist", c.Path())); os.IsNotExist(err) {
	// 		return c.SendFile("./frontend/dist/index.html", true)
	// 	}
	// 	return c.SendFile(filepath.Join("./frontend/dist", c.Path()))
	// })
	
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embeddedFiles),
		PathPrefix: "frontend/dist",
		Index:      "index.html",
		MaxAge:     3600,
	}))

	// Fallback route for SPA
	app.Get("*", func(c *fiber.Ctx) error {
        requestedPath := strings.TrimPrefix(c.Path(), "/")

        filePath := filepath.Join("frontend/dist", requestedPath)
        
        file, err := embeddedFiles.Open(filePath)
        if err != nil {
            if requestedPath != "" {
                return c.SendFile("./frontend/dist/index.html", true)
            }
            return c.SendFile("./frontend/dist/index.html")
        }
        defer file.Close()

        fileData, err := io.ReadAll(file)
        if err != nil {
            return fiber.NewError(fiber.StatusInternalServerError, "Error reading file")
        }

        return c.Send(fileData)
    })

	// Start server
	if err := app.Listen(port); err != nil {
		log.Logger.Fatal().Str("port", port).Err(err).Msg("Error starting server")
	}
	log.Logger.Info().Str("port", port).Msg("Server is listening on port")
}

func defineAPIRoutes(app *fiber.App) {
	api := app.Group("/api")
	files := api.Group("/files")
	files.Get("/", handlers.FileHandler)
	files.Post("/", handlers.FileHandler)
	files.Post("/rename", handlers.RenameHandler)
	files.Delete("/delete", handlers.DeleteHandler)
	files.Get("/view_archive", handlers.ArchiveHandler)
	files.Post("/upload", handlers.UploadFileHandler)
	files.Get("/view", handlers.ViewHandler)
	files.Post("/save", handlers.SaveHandler)
	files.Put("/permissions", handlers.UpdatePermissionsHandler)
	files.Get("/download", handlers.DownloadHandler)
	files.Get("/extract", handlers.UnzipHandler)
	files.Get("/make", handlers.MakeNewHandler)
}

func setupMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(requestLoggerMiddleware())
}

func requestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		log.Logger.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("Request")
		return err
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	log.Logger.Error().Err(err).Msg("Internal server error")
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}
