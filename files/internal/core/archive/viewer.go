package archive

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"errors"
	"files/internal/models"
	"files/internal/utils/helper"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ProcessArchiveFile(archivePath string) ([]models.ArchiveFileInfo, error) {
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		return processZipFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar.gz"), strings.HasSuffix(archivePath, ".tgz"):
		return processTarGzFile(archivePath)
	case strings.HasSuffix(archivePath, ".tar"):
		return processTarFile(archivePath)
	case strings.HasSuffix(archivePath, ".gz"):
		return processGzFile(archivePath)
	default:
		return nil, errors.New("unsupported file format")
	}
}

func processZipFile(zipPath string) ([]models.ArchiveFileInfo, error) {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipFile.Close()

	var fileInfos []models.ArchiveFileInfo
	for _, f := range zipFile.File {
		fileInfos = append(fileInfos, createArchiveFileInfo(f.Name, f.FileInfo().IsDir(), int64(f.UncompressedSize64), f.Modified))
	}
	return fileInfos, nil
}

func processTarFile(tarPath string) ([]models.ArchiveFileInfo, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TAR file: %w", err)
	}
	defer file.Close()

	return processTarReader(tar.NewReader(file))
}

func processTarGzFile(tarGzPath string) ([]models.ArchiveFileInfo, error) {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open TAR.GZ file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	return processTarReader(tar.NewReader(bufio.NewReader(gzReader)))
}

func processTarReader(tarReader *tar.Reader) ([]models.ArchiveFileInfo, error) {
	var fileInfos []models.ArchiveFileInfo
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading TAR file: %w", err)
		}
		fileInfos = append(fileInfos, createArchiveFileInfo(header.Name, header.Typeflag == tar.TypeDir, header.Size, header.ModTime))
	}
	return fileInfos, nil
}

func processGzFile(gzPath string) ([]models.ArchiveFileInfo, error) {
	file, err := os.Open(gzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open GZ file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	uncompressedSize, err := io.Copy(io.Discard, gzReader)
	if err != nil {
		return nil, fmt.Errorf("error reading GZ file: %w", err)
	}

	fileInfo := createArchiveFileInfo(filepath.Base(gzPath), false, uncompressedSize, gzReader.ModTime)
	return []models.ArchiveFileInfo{fileInfo}, nil
}

func createArchiveFileInfo(name string, isDir bool, size int64, modTime time.Time) models.ArchiveFileInfo {
	return models.ArchiveFileInfo{
		Name:         name,
		Path:         name,
		IsDir:        isDir,
		FileSize:     helper.ByteSize(size).String(),
		LastModified: modTime.Format("2006-01-02 15:04:05"),
	}
}