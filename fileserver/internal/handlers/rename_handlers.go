package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// RenameHandler handles renaming of files or directories
func RenameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	oldFileName := r.FormValue("file")
	newFileName := r.FormValue("newName")

	if oldFileName == "" || newFileName == "" {
		http.Error(w, "Both old and new names are required", http.StatusBadRequest)
		return
	}

	// Clean the old file path
	oldFilePath := filepath.Clean(oldFileName)
	extension := filepath.Ext(oldFilePath)

	// Append extension if none is provided
	if filepath.Ext(newFileName) == "" {
		newFileName += extension
	}

	newFilePath := filepath.Join(filepath.Dir(oldFilePath), newFileName)

	// Check for name conflicts
	if _, err := os.Stat(newFilePath); err == nil {
		http.Error(w, "A file with the new name already exists", http.StatusConflict)
		return
	}

	// Sanitize newFileName to prevent invalid characters
	if strings.ContainsAny(newFileName, "\\/:*?\"<>|") {
		http.Error(w, "Invalid characters in file name", http.StatusBadRequest)
		return
	}

	// Rename the file or directory
	err := os.Rename(oldFilePath, newFilePath)
	if err != nil {
		http.Error(w, "Failed to rename file", http.StatusInternalServerError)
		return
	}

	log.Printf("Renamed '%s' to '%s'", oldFilePath, newFilePath)
	http.Redirect(w, r, "/list?dir="+url.QueryEscape(filepath.Dir(oldFilePath))+"&msg=File renamed successfully", http.StatusSeeOther)
}
