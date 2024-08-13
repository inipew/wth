package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MakeNewHandler handles requests to create a new file or directory based on query parameters.
func MakeNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract parameters from query
	creationType := r.URL.Query().Get("type")
	currentPath := r.URL.Query().Get("currentPath")
	name := r.URL.Query().Get("name")

	// Validate parameters
	if err := validateParams(creationType, name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Resolve and validate the current path
	currentPath = resolvePath(currentPath)
	baseDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Failed to get the base directory", http.StatusInternalServerError)
		return
	}

	if !isValidPath(baseDir, currentPath) {
		http.Error(w, "Invalid current path", http.StatusBadRequest)
		return
	}

	// Create file or directory
	if err := createEntity(creationType, currentPath, name); err != nil {
		log.Printf("Error creating %s: %v", creationType, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the list view of the current path
	http.Redirect(w, r, "/list?dir="+filepath.ToSlash(currentPath), http.StatusSeeOther)
}

// validateParams checks the validity of the creation type and name.
func validateParams(creationType, name string) error {
	if creationType == "" || (creationType != "file" && creationType != "dir") {
		return fmt.Errorf("Invalid 'type' parameter. Must be 'file' or 'dir'.")
	}
	if name == "" {
		return fmt.Errorf("Name parameter must be provided.")
	}
	return nil
}

// resolvePath returns the cleaned and defaulted current path.
func resolvePath(path string) string {
	if path == "" {
		return "." // Default to the current working directory
	}
	return filepath.Clean(path)
}

// isValidPath checks if the new path is within the base directory.
func isValidPath(baseDir, newPath string) bool {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}
	absNewPath, err := filepath.Abs(newPath)
	if err != nil {
		return false
	}
	return strings.HasPrefix(absNewPath, absBaseDir)
}

// createEntity creates a file or directory based on the type and path provided.
func createEntity(creationType, currentPath, name string) error {
	newPath := filepath.Join(currentPath, name)

	switch creationType {
	case "dir":
		if err := os.Mkdir(newPath, os.ModePerm); err != nil {
			return fmt.Errorf("Failed to create new directory: %w", err)
		}
	case "file":
		file, err := os.Create(newPath)
		if err != nil {
			return fmt.Errorf("Failed to create new file: %w", err)
		}
		defer file.Close()
	default:
		return fmt.Errorf("Unknown type %s", creationType)
	}
	return nil
}
