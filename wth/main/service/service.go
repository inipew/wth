package service

import (
	"bufio"
	"fmt"
	"os"
	"wth/main/caddy"
	"wth/main/common"
	singbox "wth/main/sing-box"
	"wth/main/utils"
)

// SaveDomain saves the domain name to a file.
func SaveDomain(domain string) error {
	return os.WriteFile(common.DomainFilePath, []byte(domain), 0644)
}

// PromptForDomain prompts the user to enter a domain and saves it.
func PromptForDomain() error {
	fmt.Print("Enter the domain name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	domain := scanner.Text()
	if err := SaveDomain(domain); err != nil {
		return fmt.Errorf("failed to save domain: %w", err)
	}
	return nil
}

// createDirectories checks for the existence of directories and creates them if they do not exist.
func createDirectories(dirs []string) {
	for _, dir := range dirs {
		if utils.PathExists(dir) {
			fmt.Printf("Directory %s already exists.\n", dir)
		} else {
			fmt.Printf("Directory %s does not exist. Creating...\n", dir)
			if err := utils.CreateDir(dir); err != nil {
				fmt.Printf("Failed to create directory %s: %v\n", dir, err)
			} else {
				fmt.Printf("Successfully created directory %s.\n", dir)
			}
		}
	}
}

// InstallCaddy installs Caddy by calling the InstallCaddy method from the caddy package.
func InstallCaddy() error {
	return caddy.InstallCaddy()
}

// InstallSingBox installs Sing-Box by calling the InstallSingBox method from the singbox package.
func InstallSingBox() error {
	return singbox.InstallSingBox()
}

// UninstallCaddy removes the Caddy installation and its service.
func UninstallCaddy() error {
	return caddy.UninstallCaddy()
}

// UninstallSingBox removes the Sing-Box installation and its service.
func UninstallSingBox() error {
	return singbox.UninstallSingBox()
}

// InstallServices handles the installation of Caddy and Sing-Box.
func InstallServices() error {
	if err := PromptForDomain(); err != nil {
		return fmt.Errorf("error prompting for domain: %w", err)
	}

	if err := InstallCaddy(); err != nil {
		return fmt.Errorf("error installing Caddy: %w", err)
	}

	if err := InstallSingBox(); err != nil {
		return fmt.Errorf("error installing Sing-Box: %w", err)
	}

	// Define directories to check
	directories := []string{
		common.WorkDir,
		common.TmpDir,
		common.LogDir,
		common.BinDir,
		common.BackupDir,
	}

	// Create missing directories
	createDirectories(directories)

	fmt.Println("Caddy and Sing-Box installed successfully!")
	return nil
}

// UninstallServices handles the uninstallation of Caddy and Sing-Box.
func UninstallServices() error {
	if err := UninstallCaddy(); err != nil {
		return fmt.Errorf("error uninstalling Caddy: %w", err)
	}

	if err := UninstallSingBox(); err != nil {
		return fmt.Errorf("error uninstalling Sing-Box: %w", err)
	}

	fmt.Println("Caddy and Sing-Box uninstalled successfully!")
	return nil
}
