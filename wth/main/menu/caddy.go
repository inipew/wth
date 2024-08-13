package menu

import (
	"fmt"
	"wth/main/caddy"
)

// ShowCaddyMenu displays the menu options for Caddy management.
func ShowCaddyMenu() {
	for {
		fmt.Println("Caddy Management Menu:")
		fmt.Println("1. Install Caddy")
		fmt.Println("2. Start Caddy")
		fmt.Println("3. Stop Caddy")
		fmt.Println("4. Restart Caddy")
		fmt.Println("5. Change Domain in Caddyfile")
		fmt.Println("6. View Caddy Logs")
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
			caddy.InstallCaddy()
		case 2:
			caddy.StartCaddy()
		case 3:
			caddy.StopCaddy()
		case 4:
			caddy.RestartCaddy()
		case 5:
			fmt.Print("Enter new domain: ")
			var domain string
			fmt.Scan(&domain)
			caddy.ChangeDomain(domain)
		case 6:
			fmt.Println("Viewing Caddy logs is not implemented yet.")
		case 0:
			return
		default:
			fmt.Println("Invalid option, please choose again.")
		}
	}
}
