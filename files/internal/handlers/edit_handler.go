package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type FileContentRequest struct {
	FileName	string `json:"fileName"`
	Content		string `json:"content"`
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := "/"+vars["filepath"]
	if fileName == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		log.Printf("filepath: %s", fileName)
		return
	}
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
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read files", http.StatusInternalServerError)
		return
	}

	// Respond with file content
	response := FileContentRequest{
		FileName: filePath,
		Content:  string(content),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	var req FileContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(req.FileName)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Write the content back to the file
	err = os.WriteFile(absPath, []byte(req.Content), 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}