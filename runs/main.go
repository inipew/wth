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
	"github.com/rs/zerolog/log"
)

const (
	configFilePath         = "./config.ini"
	defaultPort            = ":5678"
	reloadDebounceDuration = 2 * time.Second
)

var (
	reloadTimer      *time.Timer
	reloadTimerMutex sync.Mutex
	wg               sync.WaitGroup
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS
func main() {
	logger.InitLogger()

	if err := config.LoadConfig(configFilePath); err != nil {
		log.Logger.Fatal().Err(err).Msg("Error loading config")
	}

	app := fiber.New(fiber.Config{
		DisablePreParseMultipartForm: true,
		StreamRequestBody:            true,
	})

	setupMiddleware(app)
	setupRoutes(app)

	go startServer(app)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := setupFileWatcher(ctx); err != nil {
			log.Logger.Fatal().Err(err).Msg("Error setting up file watcher")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	cancel()
	wg.Wait()
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

func setupRoutes(app *fiber.App) {
	app.Get("/api/command/list", handlers.GetCommandList)
	app.Post("/api/command/execute", handlers.CommandHandler)
	// app.Static("/", "./frontend/dist", fiber.Static{
	// 	Compress:      true,
	// 	CacheDuration: 3 * time.Hour,
	// })

	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	filePath := filepath.Join("./frontend/dist", c.Path())
	// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
	// 		return c.SendFile("./frontend/dist/index.html", true)
	// 	}
	// 	return c.SendFile(filePath)
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
}

func startServer(app *fiber.App) {
	port := config.ConfigData.WebConf.Port
	if port == "" {
		port = defaultPort
	}
	if err := app.Listen(port); err != nil {
		log.Logger.Fatal().Str("port", port).Err(err).Msg("Error starting server")
	}
	log.Logger.Info().Str("port", port).Msg("Server is listening on port")
}

func setupFileWatcher(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating file watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(configFilePath); err != nil {
		log.Logger.Error().Err(err).Msg("Error adding file to watcher")
		return fmt.Errorf("adding file to watcher: %w", err)
	}

	log.Logger.Info().Msg("Watching for changes in config file...")

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				log.Logger.Info().Msg("File watcher stopped due to context cancellation")
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Logger.Info().Msg("Config file modified; scheduling reload...")
					debounceConfigReload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Logger.Error().Err(err).Msg("Watcher error")
			}
		}
	}()

	<-done
	return nil
}

func debounceConfigReload() {
	reloadTimerMutex.Lock()
	defer reloadTimerMutex.Unlock()

	if reloadTimer != nil {
		reloadTimer.Stop()
	}
	reloadTimer = time.AfterFunc(reloadDebounceDuration, func() {
		log.Logger.Info().Msg("Reloading config...")
		if err := config.LoadConfig(configFilePath); err != nil {
			log.Logger.Warn().Err(err).Msg("Error reloading config")
		} else {
			log.Logger.Info().Str("file", configFilePath).Msg("Config file reloaded")
		}
	})
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
