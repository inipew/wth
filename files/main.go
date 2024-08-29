package main

import (
	"files/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
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
    r := mux.NewRouter()
    // Define API routes
	defineAPIRoutes(r)

	// Serve static files
	r.PathPrefix("/").HandlerFunc(staticFileHandler)

    loggedMux := LoggingMiddleware(r)

    corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change this to a more restrictive list if needed
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(loggedMux)

    server := &http.Server{
		Addr:    port,
		Handler: corsHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

    // Start server
    log.Printf("Server running on http://localhost%s\n", port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func defineAPIRoutes(r *mux.Router) {
	r.HandleFunc("/api/files", handlers.FileHandler).Methods("GET", "POST")
	r.HandleFunc("/api/files/rename", handlers.RenameHandler).Methods("POST")
	r.HandleFunc("/api/files/delete", handlers.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/api/files/view_archive", handlers.ArchiveHandler).Methods("GET")
	r.HandleFunc("/api/files/upload", handlers.UploadFileHandler).Methods("POST")
	r.HandleFunc("/api/files/view/{filepath:.*}", handlers.ViewHandler).Methods("GET")
	r.HandleFunc("/api/files/save", handlers.SaveHandler).Methods("POST")
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("./frontend/dist", r.URL.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.ServeFile(w, r, "./frontend/dist/index.html")
		return
	}
	http.FileServer(http.Dir("./frontend/dist")).ServeHTTP(w, r)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Proses permintaan
        next.ServeHTTP(w, r)
        // Log informasi tentang permintaan
        log.Printf("Method: %s, Path: %s, Duration: %s", r.Method, r.URL.Path, time.Since(start))
    })
}