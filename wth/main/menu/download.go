package menu

import (
	"context"
	"fmt"
	"log"
	"wth/main/common"
	"wth/main/utils"
)

// ShowDownloadMenu displays the menu options for downloading and updating.
func ShowDownloadMenu() {
	for {
		fmt.Println("Download Management Menu:")
		fmt.Println("1. Download Caddy")
		fmt.Println("2. Update Caddy")
		fmt.Println("3. Download Sing-Box")
		fmt.Println("4. Update Sing-Box")
		fmt.Println("5. Check Latest Version of Caddy")
		fmt.Println("6. Check Latest Version of Sing-Box")
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
			ctx := context.Background()
			// Contoh penggunaan fungsi
			err := utils.DownloadLatestRelease(ctx, "caddyserver", "caddy", common.TmpDir)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		case 2:
			fmt.Println("Updating Caddy...")
			// Call appropriate functions to update Caddy
		case 3:
			ctx := context.Background()
			// Contoh penggunaan fungsi
			err := utils.DownloadLatestRelease(ctx, "SagerNet", "sing-box", common.TmpDir)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		case 4:
			fmt.Println("Updating Caddy...")
			// Call appropriate functions to update Caddy
		case 5:
			version, err := utils.GetLatestReleaseVersion("caddyserver", "caddy")
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Printf("Latest version: %s\n", version)
		case 6:
			version, err := utils.GetLatestReleaseVersion("sagernet", "sing-box")
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Printf("Latest version: %s\n", version)
		case 0:
			return
		default:
			fmt.Println("Invalid option, please choose again.")
		}
	}
}
