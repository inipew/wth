package handlers

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"files/internal/models"
	"files/internal/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ArchiveHandler handles requests for viewing ZIP, TAR, and TAR.GZ files
func ArchiveHandler(c *fiber.Ctx) error {
	archivePath := c.Query("path")
	if archivePath == "" {
		// return c.Status(fiber.StatusBadRequest).SendString("Missing 'path' query parameter") 
		return respondWithJSON(c, fiber.StatusBadRequest,"Missing 'path' query parameter")
	}

	archivePath = filepath.Clean(archivePath)
	fileInfos, err := processArchiveFile(archivePath)
	if err != nil {
		log.Printf("Error processing archive file %s: %v", archivePath, err)
		// return c.Status(fiber.StatusInternalServerError).SendString("Failed to process archive file")
		return respondWithError(c, fiber.StatusInternalServerError,fmt.Sprintf("Error processing archive file %s: %v", archivePath, err))
	}

	return c.Status(fiber.StatusOK).JSON(models.ArchiveInfo{
		Name:  filepath.Base(archivePath),
		Path:  archivePath,
		Files: fileInfos,
	})
}

// processArchiveFile processes the archive file based on its extension
func processArchiveFile(archivePath string) ([]models.ArchiveFileInfo, error) {
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		return processZipFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar.gz"), strings.HasSuffix(archivePath, ".tgz"):
		return processAndReadTarGzFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar"):
		return processTarFile(archivePath)
	case strings.HasSuffix(archivePath, ".gz"):
		return getGzFileInfo(archivePath)
	default:
		return nil, http.ErrNotSupported
	}
}

// processZipFile processes ZIP files and returns their contents
func processZipFile(zipPath string) ([]models.ArchiveFileInfo, error) {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	var fileInfos []models.ArchiveFileInfo
	for _, f := range zipFile.File {
		fileInfos = append(fileInfos, models.ArchiveFileInfo{
			Name:          f.Name,
			Path:          f.Name,
			IsDir:         f.FileInfo().IsDir(),
			FileSize:      utils.ByteSize(int64(f.UncompressedSize64)).String(),
			LastModified:  f.Modified.Format("2006-01-02 15:04:05"),
		})
	}
	return fileInfos, nil
}

// processTarFile processes TAR files and returns their contents
func processTarFile(tarPath string) ([]models.ArchiveFileInfo, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	var fileInfos []models.ArchiveFileInfo
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fileInfos = append(fileInfos, models.ArchiveFileInfo{
			Name:          header.Name,
			Path:          header.Name,
			IsDir:         header.Typeflag == tar.TypeDir,
			FileSize:      utils.ByteSize(header.Size).String(),
			LastModified:  header.ModTime.Format("2006-01-02 15:04:05"),
		})
	}
	return fileInfos, nil
}

// processTarGzFile processes TAR.GZ files and returns their contents
func processAndReadTarGzFile(tarGzPath string) ([]models.ArchiveFileInfo, error) {
    // Open the .tar.gz file
    file, err := os.Open(tarGzPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create a gzip reader
    gzReader, err := gzip.NewReader(file)
    if err != nil {
        return nil, err
    }
    defer gzReader.Close()

    // Create a tar reader
    tarReader := tar.NewReader(bufio.NewReader(gzReader))

    var archiveFileInfos []models.ArchiveFileInfo
    for {
        header, err := tarReader.Next()
        if err != nil {
            if err == io.EOF {
                break // End of archive
            }
            return nil, err
        }

        archiveFileInfos = append(archiveFileInfos, models.ArchiveFileInfo{
            Name:          header.Name,
            Path:          header.Name,
            IsDir:         header.Typeflag == tar.TypeDir,
            FileSize:      utils.ByteSize(header.Size).String(),
            LastModified:  header.ModTime.Format("2006-01-02 15:04:05"),
        })
    }

    return archiveFileInfos, nil
}

func getGzFileInfo(gzPath string) ([]models.ArchiveFileInfo, error) {
    file, err := os.Open(gzPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    gzReader, err := gzip.NewReader(file)
    if err != nil {
        return nil, err
    }
    defer gzReader.Close()

    var uncompressedSize int64
    buf := make([]byte, 64*1024)
    for {
        n, err := gzReader.Read(buf)
        if n > 0 {
            uncompressedSize += int64(n)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
    }

    fileInfo := models.ArchiveFileInfo{
        Name:          gzReader.Name,
        Path:          gzReader.Name,
        IsDir:         false,
        FileSize:      utils.ByteSize(uncompressedSize).String(),
        LastModified:  gzReader.ModTime.Format("2006-01-02 15:04:05"),
    }

    return []models.ArchiveFileInfo{fileInfo}, nil
}