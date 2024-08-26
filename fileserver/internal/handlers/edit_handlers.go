package handlers

import (
	"encoding/json"
	"fileserver/internal/models"
	"fileserver/templates"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// type data struct {
// 	FileName string
// 	Content  string
// 	PrevDir  string
// }

// Template variable
var editTemplate *template.Template

// init function to load the template
func init() {
	var err error
	funcMap := template.FuncMap{}
	editTemplate, err = templates.GetTemplate("edit.html",funcMap)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	// fileName := r.URL.Query().Get("file")
	// if fileName == "" {
	// 	http.Error(w, "File name is required", http.StatusBadRequest)
	// 	return
	// }

	// // Decode URL encoded file path
	// decodedFileName, err := url.QueryUnescape(fileName)
	// if err != nil {
	// 	http.Error(w, "Invalid file path", http.StatusBadRequest)
	// 	return
	// }

	// // Clean the file path and resolve absolute path
	// filePath, err := filepath.Abs(filepath.Clean(decodedFileName))
	// if err != nil {
	// 	http.Error(w, "Invalid file path", http.StatusBadRequest)
	// 	return
	// }

	// // Determine the previous directory
	// prevDir := filepath.Dir(filePath)

	// log.Printf("Attempting to read file at path: %s", filePath)

	// // Read the file content
	// file, err := os.ReadFile(filePath)
	// if err != nil {
	// 	log.Printf("Error reading file: %v", err)
	// 	http.Error(w, "Failed to read file", http.StatusInternalServerError)
	// 	return
	// }

	// // Check if the file is a binary file
	// if isBinary(file) {
	// 	// Redirect to a download route instead of editing
	// 	http.Redirect(w, r, "/download?file="+url.QueryEscape(filePath), http.StatusSeeOther)
	// 	return
	// }

	// data := data{
	// 	FileName: fileName,
	// 	Content:  string(file),
	// 	PrevDir:  filepath.ToSlash(prevDir),
	// }

	if err := editTemplate.Execute(w, nil); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func isBinary(data []byte) bool {
	// Check for non-text characters
	for _, b := range data {
		if b == 0 {
			return true // Found null byte, likely binary
		}
	}
	return false
}

type FileContentRequest struct {
	FileName	string `json:"fileName"`
	Content		string `json:"content"`
	PrevDir		string `json:"prevDir"`
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	prevDir := filepath.Dir(absPath)
	// Read the file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Respond with file content
	response := FileContentRequest{
		FileName: filePath,
		Content:  string(content),
		PrevDir: prevDir,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// respondWithError mengirimkan balasan kesalahan dalam format JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.Response{Message: message})
}