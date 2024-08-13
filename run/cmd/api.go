package cmd

import (
	"context"
	"net/http"
	"os"
	"run/internal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewAPICommand creates a new command to start the API server.
func NewAPICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the API server",
		Run:   startAPI,
	}
	return cmd
}

// startAPI initializes and starts the API server.
func startAPI(cmd *cobra.Command, args []string) {
	var cfg *internal.Config
	var err error
	var cfgMutex sync.RWMutex

	// Log path and file check
	logrus.Infof("Loading config from: %s", configFile)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logrus.Fatalf("Configuration file does not exist: %s", configFile)
	}
	// Load initial config
	cfg, err = internal.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("Error loading config: %v", err)
	}

	// Create a new context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start file watcher to reload config on changes
	go func() {
		err := internal.WatchConfig(ctx, configFile, func(newCfg *internal.Config) {
			logrus.Println("Configuration file changed. Reloading...")
			cfgMutex.Lock()
			cfg = newCfg
			cfgMutex.Unlock()
		})
		if err != nil {
			logrus.Fatalf("Error watching config: %v", err)
		}
	}()

	// Define HTTP handlers
	http.HandleFunc("/api/execute", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		cfgMutex.RLock()
		internal.CommandHandler(cfg, timeout)(w, r)
		cfgMutex.RUnlock()
	}))

	http.HandleFunc("/api/list", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		cfgMutex.RLock()
		internal.ListCommandsHandler(cfg)(w, r)
		cfgMutex.RUnlock()
	}))

	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := ":5678"
	logrus.Infof("Starting API server on port %s", port)

	// Create HTTP server instance
	server := &http.Server{
		Addr:    port,
		Handler: nil,
	}

	// Run server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Failed to start API server: %v", err)
			cancel()
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	logrus.Println("Context cancelled, shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}

	logrus.Println("Server exiting")
}

// loggingMiddleware logs HTTP requests.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Received %s request for %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
