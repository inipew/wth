package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.FormValue("file")
	content := r.FormValue("content")
	prevDir := r.FormValue("prevDir")

	if fileName == "" || content == "" {
		http.Error(w, "File name and content are required", http.StatusBadRequest)
		return
	}

	// Decode and clean the file path
	decodedFileName, err := filepath.Abs(fileName)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	absFilePath, err := filepath.Abs(decodedFileName)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Ensure the file path is valid
	if err := os.WriteFile(absFilePath, []byte(content), 0644); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Redirect to the previous directory
	http.Redirect(w, r, "/list?dir="+prevDir, http.StatusSeeOther)
}
