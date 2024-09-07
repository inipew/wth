package handlers

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"files/internal/models"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func UnzipHandler(c *fiber.Ctx) error {
	filePath := c.Query("file")
	if filePath == "" {
		return respondWithError(c, fiber.StatusBadRequest, "File parameter is required")
	}

	// Decode URL encoded file path and clean it
	decodedFilePath, err := url.QueryUnescape(filePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Invalid file path")
		return respondWithError(c, fiber.StatusBadRequest, "Invalid file path: "+err.Error())
	}
	filePathClean := filepath.Clean(decodedFilePath)

	// Check if the file exists
	if _, err := os.Stat(filePathClean); os.IsNotExist(err) {
		log.Logger.Error().Err(err).Str("path", filePathClean).Msg("File does not exist")
		return respondWithError(c, fiber.StatusNotFound, "File does not exist")
	} else if err != nil {
		log.Logger.Error().Err(err).Str("path", filePathClean).Msg("Failed to access file")
		return respondWithError(c, fiber.StatusInternalServerError, "Failed to access file: "+err.Error())
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
			errMsg = "Unsupported GZ file type"
			err = extractGzipFile(filePathClean)
		}
	default:
		return respondWithError(c, fiber.StatusBadRequest, "Unsupported file type")
	}

	if err != nil {
		log.Logger.Error().Err(err).Msg(errMsg)
		return respondWithError(c, fiber.StatusBadRequest, errMsg+err.Error())
	}

	return respondWithJSON(c, fiber.StatusOK, models.Response{
		Message: "File extracted successfully",
	})
}

// unzip extracts ZIP files
func unzip(zipPath string) error {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", zipPath).Msg("Failed to open ZIP file")
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
			log.Logger.Error().Err(err).Str("path", zipPath).Msg("Error extracting ZIP file")
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

	// Create file with buffered I/O
	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to create file from ZIP entry")
		return err
	}
	defer dstFile.Close()

	// Use a buffered writer
	bufWriter := io.Writer(dstFile)

	// Copy file content using buffered I/O
	rc, err := f.Open()
	if err != nil {
		log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to open ZIP entry")
		return err
	}
	defer rc.Close()

	if _, err := io.Copy(bufWriter, rc); err != nil {
		log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to copy data from ZIP entry")
		return err
	}

	return nil
}

// untar extracts TAR files
func untar(tarPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", tarPath).Msg("Failed to open TAR file")
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
			log.Logger.Error().Err(err).Str("path", tarPath).Msg("Failed to read TAR header")
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
			log.Logger.Error().Err(err).Str("path", tarPath).Msg("Error extracting TAR file")
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

	// Create file with buffered I/O
	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to create file from TAR entry")
		return err
	}
	defer dstFile.Close()

	// Use a buffered writer
	bufWriter := io.Writer(dstFile)

	// Copy file content using buffered I/O
	if _, err := io.Copy(bufWriter, tarReader); err != nil {
		log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to copy data from TAR entry")
		return err
	}

	return nil
}

// untarGz extracts TAR.GZ files
func untarGz(tgzPath string) error {
	file, err := os.Open(tgzPath)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", tgzPath).Msg("Failed to open TAR.GZ file")
		return fmt.Errorf("failed to open TAR.GZ file: %w", err)
	}
	defer file.Close()

	// Create a GZIP reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		log.Logger.Error().Err(err).Str("path", tgzPath).Msg("Failed to create GZIP reader")
		return fmt.Errorf("failed to create GZIP reader: %w", err)
	}
	defer gzipReader.Close()

	// Create a TAR reader
	tarReader := tar.NewReader(gzipReader)

	// Determine the destination directory
	destDir := filepath.Dir(tgzPath)

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Logger.Error().Err(err).Str("directory", destDir).Msg("Failed to create destination directory")
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			log.Logger.Error().Err(err).Str("path", tgzPath).Msg("Failed to read TAR header")
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		fPath := filepath.Join(destDir, header.Name)

		// Create necessary directories
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to create directory")
				return fmt.Errorf("failed to create directory %s: %w", fPath, err)
			}
			continue
		}

		// Ensure the parent directory exists before creating the file
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to create parent directory")
			return fmt.Errorf("failed to create parent directory for %s: %w", fPath, err)
		}

		// Create file with buffered I/O
		dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to create file")
			return fmt.Errorf("failed to create file %s: %w", fPath, err)
		}

		// Copy file content using buffered I/O
		if _, err := io.Copy(dstFile, tarReader); err != nil {
			dstFile.Close()
			log.Logger.Error().Err(err).Str("path", fPath).Msg("Failed to copy data")
			return fmt.Errorf("failed to copy file content for %s: %w", fPath, err)
		}

		dstFile.Close()
	}

	return nil
}

func extractGzipFile(gzipFilePath string) error {
    // Buka file .gz
    gzFile, err := os.Open(gzipFilePath)
    if err != nil {
        log.Logger.Error().Err(err).Str("path", gzipFilePath).Msg("Error opening gzip file")
        return fmt.Errorf("error opening gzip file: %w", err)
    }
    // Pastikan file ditutup secara manual
    // defer gzFile.Close()

    // Buat reader gzip
    gzReader, err := gzip.NewReader(gzFile)
    if err != nil {
        gzFile.Close() // Pastikan file gzip ditutup jika error
        log.Logger.Error().Err(err).Str("path", gzipFilePath).Msg("Error creating gzip reader")
        return fmt.Errorf("error creating gzip reader: %w", err)
    }
    // Pastikan reader gzip ditutup secara manual
    // defer gzReader.Close()

    destFileName := filepath.Join(filepath.Dir(gzipFilePath), gzReader.Name) // Ganti dengan nama file default jika gzReader.Name kosong

    // Buka file output
    outFile, err := os.OpenFile(destFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        gzReader.Close() // Pastikan reader gzip ditutup jika error
        gzFile.Close()   // Pastikan file gzip ditutup jika error
        log.Logger.Error().Err(err).Str("path", destFileName).Msg("Error creating output file")
        return fmt.Errorf("error creating output file: %w", err)
    }
    // Pastikan file output ditutup secara manual
    // defer outFile.Close()

    // Salin data dari gzip reader ke file output
    if _, err = io.Copy(outFile, gzReader); err != nil {
        outFile.Close() // Pastikan file output ditutup jika error
        gzReader.Close() // Pastikan reader gzip ditutup jika error
        gzFile.Close()   // Pastikan file gzip ditutup jika error
        log.Logger.Error().Err(err).Str("path", destFileName).Msg("Error copying data to output file")
        return fmt.Errorf("error copying data to output file: %w", err)
    }

    // Menutup file dan reader secara eksplisit
    outFile.Close()
    gzReader.Close()
    gzFile.Close()

    return nil
}
