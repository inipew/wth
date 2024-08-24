package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// New handler for viewing the content of a file
func ViewFileHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Decode URL encoded file path
	decodedFileName, err := url.QueryUnescape(fileName)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	filePath, err := filepath.Abs(filepath.Clean(decodedFileName))
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	data := struct {
		Content string `json:"content"`
	}{
		Content: string(file),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// SaveEditHandler untuk menyimpan perubahan yang telah diedit
func SaveEditHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("file_path")
	content := r.FormValue("content")

	if filePath == "" || content == "" {
		http.Error(w, "File path and content are required", http.StatusBadRequest)
		return
	}

	// Write the new content to the file
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/api/files?path="+url.QueryEscape(filepath.Dir(filePath)), http.StatusSeeOther)
}
