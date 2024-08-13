package handlers

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fileserver/internal/utils"
	"fileserver/templates"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArchiveListing represents the structure for displaying archive contents
type ArchiveListing struct {
	ArchiveFileName string
	Files           []FileInfo
	PrevDir         string
}

// Template variable
var archiveTemplate *template.Template

// init function to load the template
func init() {
	var err error
	funcMap := template.FuncMap{}
	archiveTemplate, err = templates.GetTemplate("archive_viewer.html", funcMap)
	if err != nil {
		log.Fatalf("Failed to load template: %v", err)
	}
}

// ArchiveViewerHandler handles requests for viewing ZIP, TAR, and TAR.GZ files
func ArchiveViewerHandler(w http.ResponseWriter, r *http.Request) {
	archivePath, err := getArchivePath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	archiveFileInfos, err := processArchiveFile(archivePath)
	if err != nil {
		http.Error(w, "Failed to process archive file", http.StatusInternalServerError)
		log.Printf("Error processing archive file %s: %v", archivePath, err)
		return
	}

	// Sort files by name
	sortFileInfos(archiveFileInfos)

	// Pass data to the template
	data := ArchiveListing{
		ArchiveFileName: filepath.Base(archivePath),
		Files:           archiveFileInfos,
		PrevDir:         filepath.Dir(archivePath),
	}

	if err := renderArchiveTemplate(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error for archive_viewer.html: %v", err)
	}
}

// getArchivePath retrieves and validates the archive file path from the request
func getArchivePath(r *http.Request) (string, error) {
	archivePath := r.URL.Query().Get("file")
	if archivePath == "" {
		return "", http.ErrMissingFile
	}

	decodedArchivePath, err := url.QueryUnescape(archivePath)
	if err != nil {
		return "", err
	}
	return filepath.Clean(decodedArchivePath), nil
}

// processArchiveFile determines the type of archive and processes it accordingly
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

// renderTemplate renders the archive viewer template with the provided data
func renderArchiveTemplate(w http.ResponseWriter, data ArchiveListing) error {
	return archiveTemplate.Execute(w, data)
}

// processZipFile processes ZIP files and returns their contents
func processZipFile(zipPath string) ([]FileInfo, error) {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	return extractZipFiles(zipFile)
}

// extractZipFiles extracts file information from the provided ZIP reader
func extractZipFiles(zipFile *zip.ReadCloser) ([]FileInfo, error) {
	var archiveFileInfos []FileInfo
	for _, f := range zipFile.File {
		archiveFileInfos = append(archiveFileInfos, FileInfo{
			Name:          f.Name,
			Path:          f.Name, // Add path for further use
			IsDir:         f.FileInfo().IsDir(),
			FileSize:      int64(f.UncompressedSize64),
			FormattedSize: utils.ByteSize(int64(f.UncompressedSize64)).String(),
			LastModified:  f.Modified.Format("2006-01-02 15:04:05"),
			IsEditable:    IsFileEditable(f.Name),
		})
	}
	return archiveFileInfos, nil
}

// processTarFile processes TAR files and returns their contents
func processTarFile(tarPath string) ([]FileInfo, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readTar(bufio.NewReader(file))
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
			IsEditable:    IsFileEditable(header.Name),
		})
	}

	return archiveFileInfos, nil
}

// getGzFileInfo retrieves information from a GZ file
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
		IsEditable:    IsFileEditable(filepath.Base(gzPath)),
	}

	return []FileInfo{archiveFileInfo}, nil
}
