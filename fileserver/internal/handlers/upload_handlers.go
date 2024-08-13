package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    dir := r.FormValue("dir")
    if dir == "" {
        http.Error(w, "Directory is required", http.StatusBadRequest)
        return
    }

    // Clean and resolve the directory path
    dirPath := filepath.Clean(dir)
    absDirPath, err := filepath.Abs(dirPath)
    if err != nil {
        http.Error(w, "Invalid directory path", http.StatusBadRequest)
        return
    }

    // Ensure the directory exists
    if _, err := os.Stat(absDirPath); os.IsNotExist(err) {
        http.Error(w, "Directory does not exist", http.StatusNotFound)
        return
    }

    if err := r.ParseMultipartForm(32 << 20); // 32 MB limit
    err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file from form", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Check file extension
    // if !isValidFileExtension(header.Filename) {
    //     http.Error(w, "Invalid file type", http.StatusBadRequest)
    //     return
    // }

    // Create file in the specified directory
    filePath := filepath.Join(absDirPath, filepath.Base(header.Filename))
    out, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Failed to create file", http.StatusInternalServerError)
        return
    }
    defer out.Close()

    if _, err := io.Copy(out, file); err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    // Redirect back to the directory listing
    http.Redirect(w, r, fmt.Sprintf("/list?dir=%s", dir), http.StatusSeeOther)
}
