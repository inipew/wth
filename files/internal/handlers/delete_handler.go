package handlers

import (
	"encoding/json"
	"files/internal/utils"
	"log"
	"net/http"
	"os"
)

// DeleteHandler menangani permintaan untuk menghapus file atau direktori
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if !utils.IsValidPath(payload.Path) {
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
