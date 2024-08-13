package singbox

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"wth/main/common"
	"wth/main/utils"
)

// InstallSingBox installs Sing-Box by downloading, extracting, and setting up.
func InstallSingBox() error {
	osArch := utils.DetectOS() + "_" + utils.DetectArch()
	fmt.Printf("Detected OS and architecture: %s\n", osArch)

	version, err := utils.GetLatestReleaseVersion("sagernet", "sing-box")
	if err != nil {
		return fmt.Errorf("error getting latest Sing-Box version: %w", err)
	}
	fmt.Printf("Latest Sing-Box version: %s\n", version)

	tarFile := filepath.Join(common.TmpDir, "sing-box.tar.gz")
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

	if err := utils.ExtractTarGz(tarFile, common.TmpDir); err != nil {
		return fmt.Errorf("error extracting tar file: %w", err)
	}

	singBoxPath := filepath.Join(common.TmpDir, "sing-box")
	if err := os.Rename(singBoxPath, filepath.Join(common.BinDir, "sing-box")); err != nil {
		return fmt.Errorf("error moving Sing-Box binary: %w", err)
	}

	if err := os.Chmod(filepath.Join(common.BinDir, "sing-box"), 0755); err != nil {
		return fmt.Errorf("error setting permissions: %w", err)
	}

	if err := os.Remove(tarFile); err != nil {
		return fmt.Errorf("error removing tar file: %w", err)
	}

	serviceFile := common.SingBoxServicePath
	if err := os.WriteFile(serviceFile, []byte(common.SingBoxServiceContent), 0644); err != nil {
		return fmt.Errorf("error writing Sing-Box service file: %w", err)
	}

	fmt.Println("Sing-Box installation and configuration complete!")
	return nil
}

// StartSingBox starts the Sing-Box service.
func StartSingBox() error {
	if err := exec.Command("sudo", "systemctl", "start", "sing-box").Run(); err != nil {
		return fmt.Errorf("error starting Sing-Box: %w", err)
	}
	fmt.Println("Sing-Box started successfully!")
	return nil
}

// StopSingBox stops the Sing-Box service.
func StopSingBox() error {
	if err := exec.Command("sudo", "systemctl", "stop", "sing-box").Run(); err != nil {
		return fmt.Errorf("error stopping Sing-Box: %w", err)
	}
	fmt.Println("Sing-Box stopped successfully!")
	return nil
}

// RestartSingBox restarts the Sing-Box service.
func RestartSingBox() error {
	if err := exec.Command("sudo", "systemctl", "restart", "sing-box").Run(); err != nil {
		return fmt.Errorf("error restarting Sing-Box: %w", err)
	}
	fmt.Println("Sing-Box restarted successfully!")
	return nil
}

// UninstallSingBox removes the Sing-Box installation and its service.
func UninstallSingBox() error {
	if err := StopSingBox(); err != nil {
		return fmt.Errorf("error stopping Sing-Box: %w", err)
	}

	if err := os.Remove(common.SingBoxServicePath); err != nil {
		return fmt.Errorf("error removing Sing-Box service file: %w", err)
	}

	singBoxPath := filepath.Join(common.BinDir, "sing-box")
	if err := os.Remove(singBoxPath); err != nil {
		return fmt.Errorf("error removing Sing-Box binary: %w", err)
	}

	dirs := []string{common.SingboxDir, common.BinDir, common.BackupDir}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("Error removing directory %s: %v\n", dir, err)
		}
	}

	fmt.Println("Sing-Box uninstallation complete!")
	return nil
}
