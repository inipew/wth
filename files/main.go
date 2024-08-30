package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"files/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gopkg.in/ini.v1"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	// Load configuration
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	port := cfg.Section("server").Key("port").String()
	if port == "" {
		log.Fatal("Port not defined in config.ini")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:                  errorHandler,
		DisablePreParseMultipartForm:  true,
		StreamRequestBody:             true,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// API routes
	defineAPIRoutes(app)

	// // Serve static files
	// app.Static("/", "./frontend/dist", fiber.Static{
	// 	Compress: true,
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
	log.Printf("Server running on http://localhost%s\n", port)
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
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

func errorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}
