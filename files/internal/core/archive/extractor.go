package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func GetAndValidateFilePath(filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("file parameter is required")
	}

	decodedFilePath, err := url.QueryUnescape(filePath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	filePathClean := filepath.Clean(decodedFilePath)
	if _, err := os.Stat(filePathClean); err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("file does not exist")
		}
		return "", fmt.Errorf("failed to access file: %w", err)
	}

	return filePathClean, nil
}

func GetExtractor(filePath string) (func(string) error, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".zip":
		return unzip, nil
	case ".tar":
		return untar, nil
	case ".gz":
		if strings.HasSuffix(filePath, ".tar.gz") || strings.HasSuffix(filePath, ".tgz") {
			return untarGz, nil
		}
		return extractGzipFile, nil
	default:
		return nil, errors.New("unsupported file type")
	}
}

func unzip(zipPath string) error {
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipFile.Close()

	destDir := filepath.Dir(zipPath)
	return extractConcurrently(len(zipFile.File), func(i int, errChan chan<- error) {
		if err := extractZipFile(zipFile.File[i], filepath.Join(destDir, zipFile.File[i].Name)); err != nil {
			errChan <- err
		}
	})
}

func extractZipFile(f *zip.File, fPath string) error {
	if f.FileInfo().IsDir() {
		return os.MkdirAll(fPath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("failed to create file from ZIP entry: %w", err)
	}
	defer dstFile.Close()

	srcFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open ZIP entry: %w", err)
	}
	defer srcFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func untar(tarPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open TAR file: %w", err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	destDir := filepath.Dir(tarPath)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		if err := extractTarFile(tarReader, header, filepath.Join(destDir, header.Name)); err != nil {
			return err
		}
	}

	return nil
}

func extractTarFile(tarReader *tar.Reader, header *tar.Header, fPath string) error {
	if header.Typeflag == tar.TypeDir {
		return os.MkdirAll(fPath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("failed to create file from TAR entry: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, tarReader)
	return err
}

func untarGz(tgzPath string) error {
	file, err := os.Open(tgzPath)
	if err != nil {
		return fmt.Errorf("failed to open TAR.GZ file: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create GZIP reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	destDir := filepath.Dir(tgzPath)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		if err := extractTarFile(tarReader, header, filepath.Join(destDir, header.Name)); err != nil {
			return err
		}
	}

	return nil
}

func extractGzipFile(gzipFilePath string) error {
	gzFile, err := os.Open(gzipFilePath)
	if err != nil {
		return fmt.Errorf("error opening gzip file: %w", err)
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzReader.Close()

	destFileName := filepath.Join(filepath.Dir(gzipFilePath), gzReader.Name)
	outFile, err := os.OpenFile(destFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return fmt.Errorf("error copying data to output file: %w", err)
	}

	return nil
}

func extractConcurrently(numTasks int, extractFunc func(int, chan<- error)) error {
	var wg sync.WaitGroup
	errChan := make(chan error, numTasks)

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			extractFunc(i, errChan)
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}