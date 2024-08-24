package handlers

import (
	"encoding/json"
	"files/internal/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RenameHandler menangani permintaan untuk mengganti nama file
func RenameHandler(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsValidPath(payload.OldPath) || strings.Contains(payload.NewName, "..") {
		log.Printf("Invalid path")
		respondWithError(w,http.StatusBadRequest, "Invalid path")
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
