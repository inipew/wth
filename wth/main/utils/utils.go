package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ReleaseInfo struct {
    TagName string `json:"tag_name"`
    Name    string `json:"name"`
    Body    string `json:"body"`
}
const (
    caddyFormat   = "caddy_%s_%s_%s.tar.gz"
    singBoxFormat = "sing-box-%s-%s-%s.tar.gz"
	caddyFileName   = "caddy.tar.gz"
    singBoxFileName = "sing.tar.gz"
)

// GetLatestReleaseVersion mendapatkan versi terbaru dari sebuah repositori GitHub.
func GetLatestReleaseVersion(owner, repo string) (string, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
    
    log.Printf("Fetching latest release version from: %s", url)

    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var releaseInfo struct {
        TagName string `json:"tag_name"`
    }

    if err := json.Unmarshal(body, &releaseInfo); err != nil {
        return "", fmt.Errorf("error unmarshaling JSON: %v", err)
    }

    log.Printf("Latest release version: %s", releaseInfo.TagName)
    return releaseInfo.TagName, nil
}

// DownloadLatestRelease mendownload file dari GitHub Release ke path tujuan
func DownloadLatestRelease(ctx context.Context, owner, repo, destinationPath string) error {
    // Mendapatkan versi terbaru
    version, err := GetLatestReleaseVersion(owner, repo)
    if err != nil {
        return fmt.Errorf("error getting latest release version: %v", err)
    }

    // Menentukan format penamaan berdasarkan owner dan repo
    var fileNameFormat string
    var fileName string
    switch {
    case owner == "caddyserver" && repo == "caddy":
        fileNameFormat = caddyFormat
        fileName = fmt.Sprintf(caddyFileName) // Nama file yang diinginkan untuk lokal
    case owner == "SagerNet" && repo == "sing-box":
        fileNameFormat = singBoxFormat
        fileName = fmt.Sprintf(singBoxFileName) // Nama file yang diinginkan untuk lokal
    default:
        return fmt.Errorf("unknown repository format for owner %s and repo %s", owner, repo)
    }
	var version2 = strings.TrimPrefix(version, "v")
    // Membuat nama file berdasarkan format yang diberikan
    formattedFileName := fmt.Sprintf(fileNameFormat, version2, DetectOS(), DetectArch())
    downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", owner, repo, version, formattedFileName)

    log.Printf("Downloading file from: %s", downloadURL)

    // Konfigurasi client dengan timeout
    client := &http.Client{Timeout: 30 * time.Second}

    // Membuat request dengan context
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }

    // Mendownload file
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Menentukan path file tujuan dan membuat direktori jika belum ada
    fullPath := filepath.Join(destinationPath, fileName)
    if err := CreateDir(filepath.Dir(fullPath)); err != nil {
        return fmt.Errorf("error creating directories: %v", err)
    }

    // Menyimpan file ke disk
    if err := saveFile(fullPath, resp.Body); err != nil {
        return fmt.Errorf("error saving file: %v", err)
    }

    log.Printf("File downloaded successfully: %s", fullPath)
    return nil
}

func saveFile(filePath string, data io.Reader) error {
    outFile, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("error creating file: %v", err)
    }
    defer outFile.Close()

    if _, err := io.Copy(outFile, data); err != nil {
        return fmt.Errorf("error writing to file: %v", err)
    }

    return nil
}

// CreateDir creates a directory and sets its permissions to 755.
func CreateDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	return nil
}

// DetectArch returns the architecture of the host system.
func DetectArch() string {
	return runtime.GOARCH
}

// DetectOS returns the operating system of the host system.
func DetectOS() string {
	return runtime.GOOS
}

// RemoveDir removes a directory and its contents.
func RemoveDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", dir, err)
	}
	return nil
}

// ExtractTarGz extracts a .tar.gz file to a specified directory.
func ExtractTarGz(tarGzPath, destDir string) error {
	// Open the .tar.gz file
	file, err := os.Open(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz file %s: %w", tarGzPath, err)
	}
	defer file.Close()

	// Create a new gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create a new tar reader
	tr := tar.NewReader(gzr)

	// Iterate through the files in the archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			// Create files
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}
			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
		default:
			return fmt.Errorf("unsupported tar header type %c", header.Typeflag)
		}
	}
	return nil
}

// RemoveFile deletes a file at the given path.
func RemoveFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", filepath, err)
	}
	return nil
}

// ClearFile clears the contents of a file at the given path.
func ClearFile(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %s for clearing: %w", filepath, err)
	}
	defer f.Close()

	return nil
}

// MoveDir moves a directory from the source path to the destination path.
func MoveDir(srcPath, destPath string) error {
	cmd := exec.Command("mv", srcPath, destPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to move directory from %s to %s: %w", srcPath, destPath, err)
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// RemoveDir removes a directory if it exists.
func RemoveDirIfExists(dir string) error {
	if PathExists(dir) {
		return RemoveDir(dir)
	}
	return nil
}

// checkServiceStatus checks if a service is active.
func CheckServiceStatus(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 3 {
			// Service is not active
			return false, nil
		}
		return false, err
	}
	// Service is active
	return true, nil
}