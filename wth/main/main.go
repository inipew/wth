package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wth/main/menu"
	"wth/main/service"
	"wth/main/utils"
)

// checkServiceInstallation checks if a service is installed.
func checkServiceInstallation(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return contains(string(output), serviceName), nil
}

// contains checks if a substring is present in a string.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// checkSudo checks if the program is running with sudo privileges.
func checkSudo() bool {
	return os.Geteuid() == 0
}

func main() {
	// Check if the program is running with sudo privileges
	if !checkSudo() {
		fmt.Println("This program must be run with sudo or as root.")
		os.Exit(1)
	}

	// Check service installation and status for caddy.service and sing-box.service
	caddyInstalled, err := checkServiceInstallation("caddy.service")
	if err != nil {
		fmt.Printf("Error checking if caddy.service is installed: %v\n", err)
		os.Exit(1)
	}

	singBoxInstalled, err := checkServiceInstallation("sing-box.service")
	if err != nil {
		fmt.Printf("Error checking if sing-box.service is installed: %v\n", err)
		os.Exit(1)
	}

	caddyRunning, err := utils.CheckServiceStatus("caddy.service")
	if err != nil {
		fmt.Printf("Error checking status of caddy.service: %v\n", err)
		os.Exit(1)
	}

	singBoxRunning, err := utils.CheckServiceStatus("sing-box.service")
	if err != nil {
		fmt.Printf("Error checking status of sing-box.service: %v\n", err)
		os.Exit(1)
	}

	// Display service installation and status
	fmt.Printf("caddy.service installed: %v\n", caddyInstalled)
	fmt.Printf("caddy.service running: %v\n", caddyRunning)
	fmt.Printf("sing-box.service installed: %v\n", singBoxInstalled)
	fmt.Printf("sing-box.service running: %v\n", singBoxRunning)
	fmt.Printf("OS and Arch: %v %v\n", utils.DetectOS(), utils.DetectArch())

	// Main menu loop
	for {
		fmt.Println("Main Menu:")
		fmt.Println("1. Install Services")
		fmt.Println("2. Uninstall Services")
		fmt.Println("3. Download and Update")
		fmt.Println("4. Sing-Box Management")
		fmt.Println("5. Service Management")
		fmt.Println("0. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch choice {
		case 1:
			if err := service.InstallServices(); err != nil {
				fmt.Printf("Error installing services: %v\n", err)
			}
		case 2:
			if err := service.UninstallServices(); err != nil {
				fmt.Printf("Error uninstalling services: %v\n", err)
			}
		case 3:
			menu.ShowDownloadMenu()
		case 4:
			// Call the function for Sing-Box Management
		case 5:
			// Call the function for Service Management
		case 0:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid option, please choose again.")
		}
	}
}
