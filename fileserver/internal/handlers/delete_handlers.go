package handlers

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DeleteHandler handles requests for deleting files or directories
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Clean the file path and resolve absolute path
	filePath := filepath.Clean(fileName)

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}
	dirPath := filepath.Dir(absFilePath)

	// Attempt to delete the file or directory
	if err := os.RemoveAll(absFilePath); err != nil {
		http.Error(w, "Failed to delete file or directory", http.StatusInternalServerError)
		return
	}

	// Redirect back to the file manager
	http.Redirect(w, r, "/list?dir="+url.QueryEscape(dirPath), http.StatusSeeOther)
}
