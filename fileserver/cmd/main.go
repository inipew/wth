package main

import (
	"fileserver/internal/handlers"
	"fileserver/internal/middleware"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-ini/ini"
)

// Konfigurasi server
type Config struct {
	Port      string
	UploadDir string
}

// LoadConfig loads the server configuration from an INI file
func LoadConfig(filename string) (*Config, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	port := cfg.Section("webconf").Key("port").String()
	if port == "" {
		return nil, fmt.Errorf("port must be specified in the config file")
	}

	uploadDir := cfg.Section("webconf").Key("UploadDir").String()
	if uploadDir == "" {
		return nil, fmt.Errorf("UploadDir must be specified in the config file")
	}

	return &Config{
		Port:      port,
		UploadDir: uploadDir,
	}, nil
}


func main() {
	// Define a command-line flag for the config file path
	configPath := flag.String("config", "config.ini", "Path to the configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create uploads directory if it doesn't exist
	// if err := handlers.CreateUploadDir(config.UploadDir); err != nil {
	// 	log.Fatalf("Error creating upload directory: %v", err)
	// }

	mux := http.NewServeMux()

	// Setup routing
	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"/", handlers.IndexFileManagerHandler},
		{"/list", handlers.IndexFileManagerHandler},
		{"/upload", handlers.UploadHandler},
		{"/uploadform", handlers.UploadFormHandler},
		{"/edit", handlers.EditHandler},
		{"/delete", handlers.DeleteHandler},
		{"/save", handlers.SaveHandler},
		{"/download", handlers.DownloadHandler},
		{"/zipview", handlers.ArchiveViewerHandler},
		{"/unzip", handlers.UnzipHandler},
		// {"/rename", handlers.RenameHandler},
		{"/create", handlers.MakeNewHandler},
		{"/api/files/view", handlers.ViewHandler},
		{"/api/files/rename", handlers.RenameHandlers},

	}

	for _, route := range routes {
		mux.Handle(route.path, middleware.LoggingMiddleware(route.handler))
	}

	// Start server
	fmt.Printf("Server started at %s (upload dir: %s)\n", config.Port, config.UploadDir)
	if err := http.ListenAndServe(config.Port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}