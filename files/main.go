package main

import (
	"files/internal/handlers"
	"log"
	"net/http"

	"github.com/rs/cors"
	"gopkg.in/ini.v1"
)

func main() {
    // Load configuration from config.ini
    cfg, err := ini.Load("config.ini")
    if err != nil {
        log.Fatalf("Failed to load config file: %v", err)
    }

    port := cfg.Section("server").Key("port").String()
    if port == "" {
        log.Fatal("Port not defined in config.ini")
    }
    mux := http.NewServeMux()
    // Set up routes
    mux.HandleFunc("/api/files", handlers.FileHandler)
    mux.HandleFunc("/api/files/rename", handlers.RenameHandler)
    mux.HandleFunc("/api/files/delete", handlers.DeleteHandler)
    mux.HandleFunc("/api/files/view_archive", handlers.ArchiveHandler)
    mux.HandleFunc("/api/files/upload", handlers.UploadFileHandler)
    mux.HandleFunc("/api/files/view_file", handlers.ViewFileHandler)
    mux.HandleFunc("/api/files/save_edit", handlers.SaveEditHandler)

    mux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change this to a more restrictive list if needed
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)

    server := &http.Server{
		Addr:    port,
		Handler: corsHandler,
	}
    
    // Start server
    log.Printf("Server running on http://localhost%s\n", port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
