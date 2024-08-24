package handlers

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"files/internal/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArchiveHandler handles requests for viewing ZIP, TAR, and TAR.GZ files
func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	archivePath := r.URL.Query().Get("path")
	if archivePath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	archivePath = filepath.Clean(archivePath)
	fileInfos, err := processArchiveFile(archivePath)
	if err != nil {
		http.Error(w, "Failed to process archive file", http.StatusInternalServerError)
		log.Printf("Error processing archive file %s: %v", archivePath, err)
		return
	}

	response := struct {
		Files []FileInfo `json:"files"`
	}{
		Files: fileInfos,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// processArchiveFile processes the archive file based on its extension
func processArchiveFile(archivePath string) ([]FileInfo, error) {
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		return processZipFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar.gz"), strings.HasSuffix(archivePath, ".tgz"):
		return processTarGzFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar"):
		return processTarFile(archivePath)
	case strings.HasSuffix(archivePath, ".gz"):
		return getGzFileInfo(archivePath)
	default:
		return nil, http.ErrNotSupported
	}
}

// processZipFile processes ZIP files and returns their contents
func processZipFile(zipPath string) ([]FileInfo, error) {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	var fileInfos []FileInfo
	for _, f := range zipFile.File {
		fileInfos = append(fileInfos, FileInfo{
			Name:          f.Name,
			Path:          f.Name,
			IsDir:         f.FileInfo().IsDir(),
			FileSize:      int64(f.UncompressedSize64),
			FormattedSize: utils.ByteSize(int64(f.UncompressedSize64)).String(),
			LastModified:  f.Modified.Format("2006-01-02 15:04:05"),
			IsEditable:    false, // Set this based on your requirements
		})
	}
	return fileInfos, nil
}

// processTarFile processes TAR files and returns their contents
func processTarFile(tarPath string) ([]FileInfo, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	var fileInfos []FileInfo
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fileInfos = append(fileInfos, FileInfo{
			Name:          header.Name,
			Path:          header.Name,
			IsDir:         header.Typeflag == tar.TypeDir,
			FileSize:      header.Size,
			FormattedSize: utils.ByteSize(header.Size).String(),
			LastModified:  header.ModTime.Format("2006-01-02 15:04:05"),
			IsEditable:    false, // Set this based on your requirements
		})
	}
	return fileInfos, nil
}

// processTarGzFile processes TAR.GZ files and returns their contents
func processTarGzFile(tarGzPath string) ([]FileInfo, error) {
    file, err := os.Open(tarGzPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    gzReader, err := gzip.NewReader(file)
    if err != nil {
        return nil, err
    }
    defer gzReader.Close()

    return readTar(gzReader)
}

// readTar reads the contents of a TAR archive from an io.Reader
func readTar(reader io.Reader) ([]FileInfo, error) {
    tarReader := tar.NewReader(reader)

    var archiveFileInfos []FileInfo
    for {
        header, err := tarReader.Next()
        if err != nil {
            if err == io.EOF {
                break // End of archive
            }
            return nil, err
        }

        archiveFileInfos = append(archiveFileInfos, FileInfo{
            Name:          header.Name,
            Path:          header.Name, // Add path for further use
            IsDir:         header.Typeflag == tar.TypeDir,
            FileSize:      header.Size,
            FormattedSize: utils.ByteSize(header.Size).String(),
            LastModified:  header.ModTime.Format("2006-01-02 15:04:05"),
            IsEditable:    utils.IsFileEditable(header.Name),
        })
    }

    return archiveFileInfos, nil
}

func getGzFileInfo(gzPath string) ([]FileInfo, error) {
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

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gzReader); err != nil {
		return nil, err
	}

	uncompressedSize := int64(buf.Len())
	archiveFileInfo := FileInfo{
		Name:          filepath.Base(gzPath),
		Path:          gzPath,
		IsDir:         false,
		FileSize:      uncompressedSize,
		FormattedSize: utils.ByteSize(uncompressedSize).String(),
		LastModified:  time.Now().Format("2006-01-02 15:04:05"), // Modify as needed
		IsEditable:    utils.IsFileEditable(filepath.Base(gzPath)),
	}

	return []FileInfo{archiveFileInfo}, nil
}