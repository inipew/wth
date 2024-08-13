package caddy

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wth/main/common"
	"wth/main/utils"
)

// // GetLatestCaddy fetches the latest Caddy version from the GitHub API.
// func GetLatestCaddy() (string, error) {
// 	const githubAPIURL = "https://api.github.com/repos/caddyserver/caddy/releases/latest"
// 	version, err := utils.GetLatestVersion(githubAPIURL)
// 	if err != nil {
// 		return "", fmt.Errorf("error getting latest Caddy version: %w", err)
// 	}
// 	return version, nil
// }

// // DownloadCaddy downloads the Caddy binary tarball based on the version and architecture.
// func DownloadCaddy(version, osArch, destPath string) error {
// 	// GitHub releases URL for Caddy binaries
// 	caddyURL := fmt.Sprintf("https://github.com/caddyserver/caddy/releases/download/%s/caddy_%s.tar.gz", version, osArch)
// 	return utils.DownloadFile(caddyURL, destPath)
// }

// InstallCaddy downloads and installs Caddy.
func InstallCaddy() error {
	// Create required directories
	dirs := []string{common.TmpDir, common.LogDir, common.SingboxDir, common.CaddyDir, common.BinDir, common.BackupDir}
	for _, dir := range dirs {
		if err := utils.CreateDir(dir); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Determine OS and architecture
	osArch := utils.DetectOS() + "_" + utils.DetectArch()
	fmt.Printf("Detected OS and architecture: %s\n", osArch)

	// Fetch the latest Caddy version
	version, err := utils.GetLatestReleaseVersion("caddyserver", "caddy")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Latest version: %s\n", version)

	// Download the Caddy tar.gz file
	tarFile := filepath.Join(common.TmpDir, "caddy.tar.gz")
	ctx := context.Background()
	// Contoh penggunaan fungsi
	err = utils.DownloadLatestRelease(ctx, "caddyserver", "caddy", tarFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// Extract the tar.gz file
	if err := utils.ExtractTarGz(tarFile, common.TmpDir); err != nil {
		return fmt.Errorf("error extracting Caddy tar.gz file: %w", err)
	}

	// Move Caddy binary to the bin directory
	caddyPath := filepath.Join(common.TmpDir, "caddy")
	if err := utils.MoveDir(caddyPath, filepath.Join(common.BinDir, "caddy")); err != nil {
		return fmt.Errorf("error moving Caddy binary: %w", err)
	}

	// Set file permissions
	caddyBinaryPath := filepath.Join(common.BinDir, "caddy")
	if err := os.Chmod(caddyBinaryPath, 0755); err != nil {
		return fmt.Errorf("error setting permissions on Caddy binary: %w", err)
	}

	// Clean up
	if err := utils.RemoveDir(common.TmpDir); err != nil {
		return fmt.Errorf("error removing tmp directory: %w", err)
	}

	// Verify Caddy installation
	out, err := exec.Command(caddyBinaryPath, "version").Output()
	if err != nil {
		return fmt.Errorf("error checking Caddy version: %w", err)
	}
	if !strings.Contains(string(out), version) {
		return fmt.Errorf("installed Caddy version does not match the expected version")
	}

	// Create Caddy service file
	if err := os.WriteFile(common.CaddyServicePath, []byte(common.CaddyServiceContent), 0644); err != nil {
		return fmt.Errorf("error writing caddy.service file: %w", err)
	}

	// Create Caddyfile
	if err := os.WriteFile(common.CaddyFilePath, []byte(common.CaddyFileContent), 0644); err != nil {
		return fmt.Errorf("error writing Caddyfile: %w", err)
	}

	fmt.Println("Caddy installation and configuration complete!")
	return nil
}

// StartCaddy starts the Caddy service.
func StartCaddy() error {
	if err := exec.Command("sudo", "systemctl", "start", "caddy").Run(); err != nil {
		return fmt.Errorf("error starting Caddy: %w", err)
	}
	fmt.Println("Caddy started successfully!")
	return nil
}

// StopCaddy stops the Caddy service.
func StopCaddy() error {
	if err := exec.Command("sudo", "systemctl", "stop", "caddy").Run(); err != nil {
		return fmt.Errorf("error stopping Caddy: %w", err)
	}
	fmt.Println("Caddy stopped successfully!")
	return nil
}

// RestartCaddy restarts the Caddy service.
func RestartCaddy() error {
	if err := exec.Command("sudo", "systemctl", "restart", "caddy").Run(); err != nil {
		return fmt.Errorf("error restarting Caddy: %w", err)
	}
	fmt.Println("Caddy restarted successfully!")
	return nil
}

// ChangeDomain updates the domain in the Caddyfile and reloads Caddy.
func ChangeDomain(newDomain string) error {
	caddyFilePath := common.CaddyFilePath
	content, err := os.ReadFile(caddyFilePath)
	if err != nil {
		return fmt.Errorf("error reading Caddyfile: %w", err)
	}

	// Replace the domain in the file content
	newContent := strings.Replace(string(content), "example.com", newDomain, 1)
	if err := os.WriteFile(caddyFilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("error writing to Caddyfile: %w", err)
	}

	// Reload Caddy to apply the changes
	if err := exec.Command("sudo", "systemctl", "reload", "caddy").Run(); err != nil {
		return fmt.Errorf("error reloading Caddy: %w", err)
	}

	fmt.Println("Domain in Caddyfile updated successfully and Caddy reloaded!")
	return nil
}

// UninstallCaddy removes the Caddy installation and its service.
func UninstallCaddy() error {
	// Stop the Caddy service if it's running
	if running, err := utils.CheckServiceStatus("caddy.service"); err != nil {
		return fmt.Errorf("error checking status of caddy.service: %w", err)
	} else if running {
		if err := exec.Command("sudo", "systemctl", "stop", "caddy").Run(); err != nil {
			return fmt.Errorf("error stopping Caddy service: %w", err)
		}
		fmt.Println("Caddy service stopped.")
	}

	// Remove the Caddy service file
	if err := os.Remove("/etc/systemd/system/caddy.service"); err != nil {
		return fmt.Errorf("error removing Caddy service file: %w", err)
	}
	fmt.Println("Caddy service file removed.")

	// Remove the Caddy binary
	if err := os.Remove(filepath.Join(common.BinDir, "caddy")); err != nil {
		return fmt.Errorf("error removing Caddy binary: %w", err)
	}
	fmt.Println("Caddy binary removed.")

	// Remove directories if empty
	dirs := []string{common.CaddyDir, common.BackupDir}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("error removing directory %s: %w", dir, err)
		}
		fmt.Printf("Directory %s removed.\n", dir)
	}

	fmt.Println("Caddy uninstallation complete!")
	return nil
}