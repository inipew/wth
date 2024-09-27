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
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

var log *logger.Logger

func main() {
	if err := run(); err != nil {
		log.Fatal("Application failed to start", "error", err)
	}
}

func run() error {
	// Initialize logger
	logConfig := logger.DefaultConfig()
	logConfig.Level = logger.Debug
	logConfig.UseJSON = false
	log = logger.New(logConfig)

	log.Info("Initializing application")

	cfg, err := config.Load("config.ini")
	if err != nil {
		return err
	}

	app := createFiberApp(cfg)

	setupMiddleware(app,log)
	handlers := handlers.NewHandlers(cfg, log)
	defineAPIRoutes(app, handlers)
	setupStaticFileServing(app, cfg.Server.UseEmbeddedFiles)

	return startServer(app, cfg.Server.Port)
}

func createFiberApp(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler:                 errorHandler,
		DisablePreParseMultipartForm: true,
		StreamRequestBody:            true,
		ReadTimeout:                  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:                 time.Duration(cfg.Server.WriteTimeout) * time.Second,
	})
}

func setupMiddleware(app *fiber.App, log *logger.Logger) {
	middleware.SetupCORS(app)
	middleware.SetupCompression(app)
	app.Use(middleware.RequestLogger(log))
}

func defineAPIRoutes(app *fiber.App, h *handlers.Handlers) {
	api := app.Group("/api")
	files := api.Group("/files")

	files.Get("/", h.FileHandler)
	files.Post("/", h.FileHandler)
	files.Post("/rename", h.RenameHandler)
	files.Delete("/delete", h.DeleteHandler)
	files.Get("/view_archive", h.ArchiveHandler)
	files.Post("/upload", h.UploadFileHandler)
	files.Get("/view", h.ViewHandler)
	files.Post("/save", h.SaveHandler)
	files.Put("/permissions", h.UpdatePermissionsHandler)
	files.Get("/download", h.DownloadHandler)
	files.Get("/extract", h.ExtractorHandler)
	files.Get("/make", h.MakeNewHandler)
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
	log.Info("Server is listening on", "port", port)
	return app.Listen(port)
}

func errorHandler(c *fiber.Ctx, err error) error {
	log.Error("Internal server error", err)
	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
}