package handlers

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func UnzipHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File parameter is required", http.StatusBadRequest)
		return
	}

	// Decode URL encoded file path and clean it
	decodedFilePath, err := url.QueryUnescape(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}
	filePathClean := filepath.Clean(decodedFilePath)

	// Check if the file exists
	if _, err := os.Stat(filePathClean); os.IsNotExist(err) {
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to access file", http.StatusInternalServerError)
		return
	}

	// Determine the file extension
	ext := strings.ToLower(filepath.Ext(filePathClean))

	var errMsg string

	switch ext {
	case ".zip":
		errMsg = "Failed to extract ZIP file: "
		err = unzip(filePathClean)
	case ".tar":
		errMsg = "Failed to extract TAR file: "
		err = untar(filePathClean)
	case ".gz":
		// Check if the file is a TAR.GZ file by looking for .tar in the name
		if strings.HasSuffix(filePathClean, ".tar.gz") || strings.HasSuffix(filePathClean, ".tgz") {
			errMsg = "Failed to extract TAR.GZ file: "
			err = untarGz(filePathClean)
		} else {
			http.Error(w, "Unsupported GZ file type", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, errMsg+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the file manager with a success message
	destDir := filepath.Dir(filePathClean)
	http.Redirect(w, r, "/list?dir="+url.QueryEscape(destDir), http.StatusSeeOther)
}


// unzip extracts ZIP files
func unzip(zipPath string) error {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Determine the destination directory
	destDir := filepath.Dir(zipPath)

	// Use a wait group to handle concurrency
	var wg sync.WaitGroup
	errChan := make(chan error, len(zipFile.File))

	// Extract files concurrently
	for _, f := range zipFile.File {
		wg.Add(1)
		go func(f *zip.File) {
			defer wg.Done()
			fPath := filepath.Join(destDir, f.Name)
			if err := extractZipFile(f, fPath); err != nil {
				errChan <- err
			}
		}(f)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a ZIP archive
func extractZipFile(f *zip.File, fPath string) error {
	if f.FileInfo().IsDir() {
		// Create directories
		return os.MkdirAll(fPath, os.ModePerm)
	}

	// Ensure the parent directory exists before creating the file
	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		return err
	}

	// Check if the file already exists
	if _, err := os.Stat(fPath); err == nil {
		// Confirm overwriting
		if err := confirmOverwrite(fPath); err != nil {
			return err
		}
	}

	// Create file with buffered I/O
	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Use a buffered writer
	bufWriter := io.Writer(dstFile)

	// Copy file content using buffered I/O
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(bufWriter, rc); err != nil {
		return err
	}

	return nil
}

// untar extracts TAR files
func untar(tarPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	// Determine the destination directory
	destDir := filepath.Dir(tarPath)

	// Use a wait group to handle concurrency
	var wg sync.WaitGroup
	errChan := make(chan error)

	// Extract files concurrently
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(header *tar.Header) {
			defer wg.Done()
			fPath := filepath.Join(destDir, header.Name)
			if err := extractTarFile(tarReader, header, fPath); err != nil {
				errChan <- err
			}
		}(header)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarFile extracts a single file from a TAR archive
func extractTarFile(tarReader *tar.Reader, header *tar.Header, fPath string) error {
	if header.Typeflag == tar.TypeDir {
		// Create directories
		return os.MkdirAll(fPath, os.ModePerm)
	}

	// Ensure the parent directory exists before creating the file
	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		return err
	}

	// Check if the file already exists
	if _, err := os.Stat(fPath); err == nil {
		// Confirm overwriting
		if err := confirmOverwrite(fPath); err != nil {
			return err
		}
	}

	// Create file with buffered I/O
	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Use a buffered writer
	bufWriter := io.Writer(dstFile)

	// Copy file content using buffered I/O
	if _, err := io.Copy(bufWriter, tarReader); err != nil {
		return err
	}

	return nil
}

// untarGz extracts TAR.GZ files
func untarGz(tgzPath string) error {
	file, err := os.Open(tgzPath)
	if err != nil {
		return fmt.Errorf("failed to open TAR.GZ file: %w", err)
	}
	defer file.Close()

	// Create a GZIP reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create GZIP reader: %w", err)
	}
	defer gzipReader.Close()

	// Create a TAR reader
	tarReader := tar.NewReader(gzipReader)

	// Determine the destination directory
	destDir := filepath.Dir(tgzPath)

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		fPath := filepath.Join(destDir, header.Name)

		// Create necessary directories
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fPath, err)
			}
			continue
		}

		// Ensure the parent directory exists before creating the file
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", fPath, err)
		}

		// Create file with buffered I/O
		dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", fPath, err)
		}

		// Copy file content using buffered I/O
		if _, err := io.Copy(dstFile, tarReader); err != nil {
			dstFile.Close()
			return fmt.Errorf("failed to copy file content for %s: %w", fPath, err)
		}

		dstFile.Close()
	}

	return nil
}

// confirmOverwrite prompts the user to confirm overwriting an existing file
func confirmOverwrite(filePath string) error {
	// For this example, we'll just return an error.
	// In a real application, you might return a message to the user.
	return fmt.Errorf("File %s already exists. Overwrite? (y/n)", filePath)
}
