package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sbx/internal/fileutils"

	"github.com/rs/zerolog/log"
)

// UntarGz extracts a TAR.GZ file to the specified destination directory, skipping directories.
func UntarGz(tgzPath, destDir string) error {
	file, err := os.Open(tgzPath)
	if err != nil {
		return fmt.Errorf("failed to open TAR.GZ file %s: %w", tgzPath, err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create GZIP reader from file %s: %w", tgzPath, err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	if err := fileutils.CreateDir(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	if err := extractFiles(tarReader, destDir); err != nil {
		return err
	}

	log.Info().Str("source", tgzPath).Str("destination", destDir).Msg("Successfully extracted archive")
	return nil
}

// extractFiles extracts files from the TAR reader to the destination directory,
// ignoring the directory structure in the archive, and skipping README.md and LICENSE.
func extractFiles(tarReader *tar.Reader, destDir string) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of tar archive
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue // Skip non-regular files
		}

		baseName := filepath.Base(header.Name)
		if baseName == "README.md" || baseName == "LICENSE" {
			log.Info().Msgf("Skipping file: %s", header.Name)
			continue
		}

		destPath := filepath.Join(destDir, baseName)
		if err := createFile(destPath, header, tarReader); err != nil {
			return err
		}
	}
	return nil
}

// createFile creates a file for the TAR entry and copies the content.
func createFile(fPath string, header *tar.Header, tarReader *tar.Reader) error {
	if err := fileutils.CreateDir(filepath.Dir(fPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", fPath, err)
	}

	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fPath, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, tarReader); err != nil {
		return fmt.Errorf("failed to copy content to %s: %w", fPath, err)
	}

	log.Info().Msgf("Successfully extracted file: %s", fPath)
	return nil
}