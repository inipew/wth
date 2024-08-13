package menu

import (
	"fmt"
	singbox "wth/main/sing-box"
)

// ShowSingBoxMenu displays the menu options for Sing-Box management.
func ShowSingBoxMenu() {
	for {
		fmt.Println("Sing-Box Management Menu:")
		fmt.Println("1. Install Sing-Box")
		fmt.Println("2. Start Sing-Box")
		fmt.Println("3. Stop Sing-Box")
		fmt.Println("4. Restart Sing-Box")
		fmt.Println("5. Change Domain in Sing-Box Config")
		fmt.Println("6. View Sing-Box Logs")
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
			singbox.InstallSingBox()
		case 2:
			singbox.StartSingBox()
		case 3:
			singbox.StopSingBox()
		case 4:
			singbox.RestartSingBox()
		case 5:
			fmt.Println("Viewing Sing-Box logs is not implemented yet.")
		case 0:
			return
		default:
			fmt.Println("Invalid option, please choose again.")
		}
	}
}
