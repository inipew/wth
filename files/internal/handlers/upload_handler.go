package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get the uploaded file
		file, header, err := r.FormFile("file") // Ambil file dan header
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read the file contents
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}

		// Determine the target path
		path := r.FormValue("path")
		filePath := filepath.Join(path, header.Filename) // Gunakan header.Filename

		// Create a new file
		err = os.WriteFile(filePath, fileBytes, 0644)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
