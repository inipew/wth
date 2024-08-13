package handlers

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CompressHandler handles requests to compress files or folders
func CompressHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file or directory paths from the query parameters
	filesParam := r.URL.Query().Get("files")
	if filesParam == "" {
		http.Error(w, "Files parameter is required", http.StatusBadRequest)
		return
	}

	files := strings.Split(filesParam, ",")
	archiveType := r.URL.Query().Get("type") // "zip", "tar", or "tar.gz"
	archiveName := r.URL.Query().Get("name")  // Archive name (without extension)

	if archiveName == "" {
		http.Error(w, "Archive name parameter is required", http.StatusBadRequest)
		return
	}

	var err error
	switch archiveType {
	case "zip":
		err = createZipArchive(w, files, archiveName)
	case "tar":
		err = createTarArchive(w, files, archiveName)
	case "tar.gz":
		err = createTarGzArchive(w, files, archiveName)
	default:
		http.Error(w, "Unsupported archive type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to create archive", http.StatusInternalServerError)
		log.Printf("Error creating archive: %v", err)
		return
	}
}

// createZipArchive creates a ZIP archive from the specified files and writes it to the response writer
func createZipArchive(w http.ResponseWriter, files []string, archiveName string) error {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", archiveName))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

// addFileToZip adds a file to the ZIP archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}

// createTarArchive creates a TAR archive from the specified files and writes it to the response writer
func createTarArchive(w http.ResponseWriter, files []string, archiveName string) error {
	w.Header().Set("Content-Type", "application/x-tar")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar", archiveName))

	tarWriter := tar.NewWriter(w)
	defer tarWriter.Close()

	for _, file := range files {
		if err := addFileToTar(tarWriter, file); err != nil {
			return err
		}
	}
	return nil
}

// addFileToTar adds a file to the TAR archive
func addFileToTar(tarWriter *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath) // set the file name in the TAR header
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}

// createTarGzArchive creates a TAR.GZ archive from the specified files and writes it to the response writer
func createTarGzArchive(w http.ResponseWriter, files []string, archiveName string) error {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for _, file := range files {
		if err := addFileToTar(tarWriter, file); err != nil {
			return err
		}
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}
	if err := gzWriter.Close(); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar.gz", archiveName))

	_, err := w.Write(buf.Bytes())
	return err
}
