package menu

import (
	"fmt"
	"os/exec"
)

// ShowServiceMenu displays the menu options for service management.
func ShowServiceMenu() {
	for {
		fmt.Println("Service Management Menu:")
		fmt.Println("1. Install Caddy and Sing-Box Services")
		fmt.Println("2. Uninstall Caddy and Sing-Box")
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
			fmt.Println("Installing Caddy and Sing-Box services...")
			// Implement installation logic
		case 2:
			fmt.Println("Uninstalling Caddy and Sing-Box...")
			err := exec.Command("sudo", "systemctl", "stop", "caddy").Run()
			if err != nil {
				fmt.Println("Error stopping Caddy:", err)
			}
			err = exec.Command("sudo", "systemctl", "disable", "caddy").Run()
			if err != nil {
				fmt.Println("Error disabling Caddy:", err)
			}
			err = exec.Command("sudo", "rm", "-f", "/etc/systemd/system/caddy.service").Run()
			if err != nil {
				fmt.Println("Error removing Caddy service file:", err)
			}
			err = exec.Command("sudo", "systemctl", "stop", "sing-box").Run()
			if err != nil {
				fmt.Println("Error stopping Sing-Box:", err)
			}
			err = exec.Command("sudo", "systemctl", "disable", "sing-box").Run()
			if err != nil {
				fmt.Println("Error disabling Sing-Box:", err)
			}
			err = exec.Command("sudo", "rm", "-f", "/etc/systemd/system/sing-box.service").Run()
			if err != nil {
				fmt.Println("Error removing Sing-Box service file:", err)
			}
		case 0:
			return
		default:
			fmt.Println("Invalid option, please choose again.")
		}
	}
}
