package main

import (
	"embed"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"files/internal/api/handlers"
	"files/internal/config"
	"files/internal/middleware"
	"files/internal/utils/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rs/zerolog/log"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Application failed to start")
	}
}

func run() error {
	if err := logger.InitLogger(); err != nil {
		return err
	}

	cfg, err := config.Load("config.ini")
	if err != nil {
		return err
	}

	app := createFiberApp()

	setupMiddleware(app)
	defineAPIRoutes(app)
	setupStaticFileServing(app, cfg.Server.UseEmbeddedFiles)

	return startServer(app, cfg.Server.Port)
}

func createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler:                 errorHandler,
		DisablePreParseMultipartForm: true,
		StreamRequestBody:            true,
	})
}

func setupMiddleware(app *fiber.App) {
	middleware.SetupCORS(app)
	middleware.SetupCompression(app)
	app.Use(middleware.RequestLogger())
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

func setupStaticFileServing(app *fiber.App, useEmbeddedFiles bool) {
	if useEmbeddedFiles {
		app.Use("/", filesystem.New(filesystem.Config{
			Root:       http.FS(embeddedFiles),
			PathPrefix: "frontend/dist",
			Index:      "index.html",
			MaxAge:     3600,
		}))

		app.Get("*", serveEmbeddedSPA)
	} else {
		app.Static("/", "./frontend/dist", fiber.Static{
			Compress:      true,
			CacheDuration: 3 * time.Hour,
		})

		app.Get("/*", serveDevelopmentSPA)
	}
}

func serveEmbeddedSPA(c *fiber.Ctx) error {
	requestedPath := strings.TrimPrefix(c.Path(), "/")
	filePath := filepath.Join("frontend/dist", requestedPath)

	file, err := embeddedFiles.Open(filePath)
	if err != nil {
		return c.SendFile("frontend/dist/index.html")
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading file")
	}

	return c.Send(fileData)
}

func serveDevelopmentSPA(c *fiber.Ctx) error {
	if _, err := os.Stat(filepath.Join("./frontend/dist", c.Path())); os.IsNotExist(err) {
		return c.SendFile("./frontend/dist/index.html", true)
	}
	return c.SendFile(filepath.Join("./frontend/dist", c.Path()))
}

func startServer(app *fiber.App, port string) error {
	log.Info().Str("port", port).Msg("Server is listening on port")
	return app.Listen(port)
}

func errorHandler(c *fiber.Ctx, err error) error {
	log.Error().Err(err).Msg("Internal server error")
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}