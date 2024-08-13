package caddy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/inipew/wth/caddy/config" // Ganti dengan path modul Anda jika diperlukan
)

// InstallCaddy installs Caddy and sets up the necessary files.
func InstallCaddy() {
	commands := []struct {
		cmd  string
		args []string
	}{
		{"sudo", []string{"apt", "install", "-y", "debian-keyring", "debian-archive-keyring", "apt-transport-https", "curl"}},
		{"curl", []string{"-1sLf", "https://dl.cloudsmith.io/public/caddy/stable/gpg.key", "|", "sudo", "gpg", "--dearmor", "-o", "/usr/share/keyrings/caddy-stable-archive-keyring.gpg"}},
		{"curl", []string{"-1sLf", "https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt", "|", "sudo", "tee", "/etc/apt/sources.list.d/caddy-stable.list"}},
		{"sudo", []string{"apt", "update"}},
		{"sudo", []string{"apt", "install", "caddy"}},
	}

	for _, cmd := range commands {
		fmt.Printf("Running command: %s %v\n", cmd.cmd, cmd.args)
		err := exec.Command(cmd.cmd, cmd.args...).Run()
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			return
		}
	}

	// Create caddy.service file
	err := os.WriteFile("/etc/systemd/system/caddy.service", []byte(config.CaddyServiceContent), 0644)
	if err != nil {
		fmt.Printf("Error writing caddy.service file: %v\n", err)
		return
	}

	// Create Caddyfile
	err = os.WriteFile("/etc/caddy/Caddyfile", []byte(config.CaddyFileContent), 0644)
	if err != nil {
		fmt.Printf("Error writing Caddyfile: %v\n", err)
		return
	}

	fmt.Println("Caddy installation and configuration complete!")
}

// StartCaddy starts the Caddy service.
func StartCaddy() {
	err := exec.Command("sudo", "systemctl", "start", "caddy").Run()
	if err != nil {
		fmt.Printf("Error starting Caddy: %v\n", err)
		return
	}
	fmt.Println("Caddy started successfully!")
}

// RestartCaddy restarts the Caddy service.
func RestartCaddy() {
	err := exec.Command("sudo", "systemctl", "restart", "caddy").Run()
	if err != nil {
		fmt.Printf("Error restarting Caddy: %v\n", err)
		return
	}
	fmt.Println("Caddy restarted successfully!")
}

// StopCaddy stops the Caddy service.
func StopCaddy() {
	err := exec.Command("sudo", "systemctl", "stop", "caddy").Run()
	if err != nil {
		fmt.Printf("Error stopping Caddy: %v\n", err)
		return
	}
	fmt.Println("Caddy stopped successfully!")
}

// ChangeDomain updates the domain in the Caddyfile.
func ChangeDomain(newDomain string) {
	caddyFilePath := "/etc/caddy/Caddyfile"
	content, err := os.ReadFile(caddyFilePath)
	if err != nil {
		fmt.Printf("Error reading Caddyfile: %v\n", err)
		return
	}

	// Replace the domain in the file content
	newContent := strings.Replace(string(content), "example.com", newDomain, 1)
	err = os.WriteFile(caddyFilePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing to Caddyfile: %v\n", err)
		return
	}

	// Reload Caddy to apply the changes
	err = exec.Command("sudo", "systemctl", "reload", "caddy").Run()
	if err != nil {
		fmt.Printf("Error reloading Caddy: %v\n", err)
		return
	}

	fmt.Println("Domain in Caddyfile updated successfully and Caddy reloaded!")
}
