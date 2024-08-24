package main

import (
	"files/internal/handlers"
	"log"
	"net/http"

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

    // Set up routes
    http.HandleFunc("/api/files", handlers.FileHandler)
    http.HandleFunc("/api/files/rename", handlers.RenameHandler)
    http.HandleFunc("/api/files/delete", handlers.DeleteHandler)

    // Start server
    log.Printf("Server running on http://localhost%s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
