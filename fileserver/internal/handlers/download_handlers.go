package handlers

import (
	"net/http"
	"path/filepath"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
    fileParam := r.URL.Query().Get("file")
    if fileParam == "" {
        http.Error(w, "File parameter is required", http.StatusBadRequest)
        return
    }

    // Clean and get the absolute path of the file
    absFilePath, err := filepath.Abs(filepath.Clean(fileParam))
    if err != nil {
        http.Error(w, "Invalid file path", http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(absFilePath))
    w.Header().Set("Content-Type", "application/octet-stream")
    // Serve the file
    http.ServeFile(w, r, absFilePath)
}
