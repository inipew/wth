package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Response structure for API responses
type Response struct {
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"`
}

// FileHandler handles requests to list files in a directory
func FileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    dir := "." // Ganti dengan direktori target Anda jika perlu
    files, err := os.ReadDir(dir)
    if err != nil {
        log.Printf("Error reading directory: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to read directory")
        return
    }

    var fileList []string
    for _, file := range files {
        fileList = append(fileList, file.Name())
    }

    respondWithJSON(w, http.StatusOK, fileList)
}

// RenameHandler handles requests to rename a file or directory
func RenameHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var data struct {
        OldName string `json:"old_name"`
        NewName string `json:"new_name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Printf("Error decoding request body: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    if err := os.Rename(data.OldName, data.NewName); err != nil {
        log.Printf("Error renaming file: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to rename file")
        return
    }

    respondWithJSON(w, http.StatusOK, Response{Message: "File renamed successfully"})
}

// DeleteHandler handles requests to delete a file or directory
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var data struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        log.Printf("Error decoding request body: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    if err := os.Remove(data.Name); err != nil {
        log.Printf("Error deleting file: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete file")
        return
    }

    respondWithJSON(w, http.StatusOK, Response{Message: "File deleted successfully"})
}

// respondWithJSON sends a JSON response to the client
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}

// respondWithError sends an error response to the client
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, Response{Message: message})
}
