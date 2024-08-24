package handlers

import (
	"encoding/json"
	"files/internal/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxPathDepth = 10 // Maksimal kedalaman direktori yang diizinkan
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type FileInfo struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	IsDir         bool   `json:"is_dir"`
	FileSize      int64  `json:"file_size"`
	FormattedSize string `json:"formatted_size"`
	LastModified  string `json:"last_modified"`
	IsEditable    bool   `json:"is_editable"`
}

type DirectoryInfo struct {
	CurrentPath   string     `json:"current_path"`
	PreviousPath  string     `json:"previous_path,omitempty"`
	Files         []FileInfo `json:"files"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func isValidPath(path string) bool {
	// Validasi untuk memastikan path tidak mengarah ke luar direktori yang diizinkan
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return !strings.Contains(absPath, "..") && len(strings.Split(absPath, string(os.PathSeparator))) <= maxPathDepth
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentPath := r.URL.Query().Get("path")
	if currentPath == "" {
		currentPath = "."
	}

	if !isValidPath(currentPath) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	previousPath := filepath.Dir(currentPath)
	if currentPath == previousPath {
		previousPath = ""
	}

	files, err := os.ReadDir(currentPath)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to read directory")
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info: %v", err)
			continue
		}

		fileInfos = append(fileInfos, FileInfo{
			Name:          file.Name(),
			Path:          filepath.Join(currentPath, file.Name()),
			IsDir:         file.IsDir(),
			FileSize:      info.Size(),
			FormattedSize: utils.ByteSize(info.Size()).String(),
			LastModified:  info.ModTime().Format(time.RFC3339),
			IsEditable:    !file.IsDir(),
		})
	}

	respondWithJSON(w, http.StatusOK, DirectoryInfo{
		CurrentPath:  currentPath,
		PreviousPath: previousPath,
		Files:        fileInfos,
	})
}

func RenameHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		OldPath string `json:"oldPath"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !isValidPath(payload.OldPath) || strings.Contains(payload.NewName, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	newPath := filepath.Join(filepath.Dir(payload.OldPath), payload.NewName)

	if err := os.Rename(payload.OldPath, newPath); err != nil {
		log.Printf("Failed to rename file: %v", err)
		http.Error(w, "Failed to rename file", http.StatusInternalServerError)
		return
	}

	log.Println("File renamed successfully")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !isValidPath(payload.Path) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if err := os.RemoveAll(payload.Path); err != nil {
		log.Printf("Failed to delete file: %v", err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{Message: message})
}
