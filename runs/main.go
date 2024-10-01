package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runs/internal/config"
	"runs/internal/handlers"
	"runs/internal/logger"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rs/zerolog"
)

const (
	configFilePath         = "./config.ini"
	defaultPort            = ":5678"
	reloadDebounceDuration = 2 * time.Second
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

type application struct {
	config        *config.ConfigManager
	app           *fiber.App
	reloadTimer   *time.Timer
	reloadMutex   sync.Mutex
	shutdownChan  chan os.Signal
}

func main() {
	cfg := logger.DefaultConfig()
    cfg.LogLevel = zerolog.DebugLevel
    // cfg.OutputFile = "app.log"
	cfg.UseColor = true

    err := logger.InitGlobalLogger(cfg)
    if err != nil {
        panic(err)
    }

	app := newApplication()
	app.run()
}

func newApplication() *application {
	configProvider := config.NewFileConfigProvider(configFilePath, *logger.GetLogger())
	configManager := config.NewConfigManager(configProvider, *logger.GetLogger())

	if err := configManager.Load(context.Background()); err != nil {
		logger.GetLogger().Fatal().Err(err).Msg("Error loading config")
	}
	// configManager.LogConfig()

	fiberApp := fiber.New(fiber.Config{
		DisablePreParseMultipartForm: true,
		StreamRequestBody:            true,
	})

	return &application{
		config:       configManager,
		app:          fiberApp,
		shutdownChan: make(chan os.Signal, 1),
	}
}

func (a *application) run() {
	a.setupMiddleware()
	a.setupRoutes()

	go a.startServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.setupFileWatcher(ctx); err != nil {
			logger.GetLogger().Fatal().Err(err).Msg("Error setting up file watcher")
		}
	}()

	signal.Notify(a.shutdownChan, os.Interrupt)
	<-a.shutdownChan

	logger.GetLogger().Info().Msg("Shutting down gracefully...")
	cancel()
	wg.Wait()
	if err := a.app.Shutdown(); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Error during server shutdown")
	}
}

func (a *application) setupMiddleware() {
	a.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	a.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	a.app.Use(a.requestLoggerMiddleware())
}

func (a *application) setupRoutes() {
	handler := handlers.NewHandler(a.config)

	a.app.Get("/api/command/list", handler.GetCommandList)
	a.app.Post("/api/command/execute", handler.CommandHandler)

	a.app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embeddedFiles),
		PathPrefix: "frontend/dist",
		Index:      "index.html",
		MaxAge:     3600,
	}))

	a.app.Get("*", a.spaHandler)
}

func (a *application) startServer() {
	port := a.config.GetConfig().WebConf.Port
	if port == "" {
		port = defaultPort
	}
	logger.GetLogger().Info().Str("port", port).Msg("Server is listening on")
	if err := a.app.Listen(port); err != nil {
		logger.GetLogger().Fatal().Str("port", port).Err(err).Msg("Error starting server")
	}
}

func (a *application) setupFileWatcher(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating file watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(configFilePath); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Error adding file to watcher")
		return fmt.Errorf("adding file to watcher: %w", err)
	}

	logger.GetLogger().Info().Msg("Watching for changes in config file...")

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info().Msg("File watcher stopped due to context cancellation")
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.GetLogger().Info().Msg("Config file modified; scheduling reload...")
				a.debounceConfigReload()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			logger.GetLogger().Error().Err(err).Msg("Watcher error")
		}
	}
}

func (a *application) debounceConfigReload() {
	a.reloadMutex.Lock()
	defer a.reloadMutex.Unlock()

	if a.reloadTimer != nil {
		a.reloadTimer.Stop()
	}
	a.reloadTimer = time.AfterFunc(reloadDebounceDuration, func() {
		logger.GetLogger().Info().Msg("Reloading config...")
		if err := a.config.Load(context.Background()); err != nil {
			logger.GetLogger().Warn().Err(err).Msg("Error reloading config")
		} else {
			logger.GetLogger().Info().Str("file", configFilePath).Msg("Config file reloaded")
		}
	})
}

func (a *application) requestLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		logger.GetLogger().Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("Request")
		return err
	}
}

func (a *application) spaHandler(c *fiber.Ctx) error {
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
}