package main

import (
	"fmt"
	"log"
	"singconfig/internal/config"
)

func main() {
	singconfig := config.BuildSingBoxConfig()
	if err := singconfig.SaveToFile("config.json"); err != nil {
		log.Fatalf("Error saving JSON file: %s", err.Error())
	}
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Example usage:
	// Add user
	config.AddUser(cfg, "", "socks", "akupew")

	// Remove user
	// config.RemoveUser(cfg, "default")

	// Update DNS
	// config.UpdateDNS(cfg, "dns.google", "https", "prefer_ipv6")
	config.DisplayInboundDetails(cfg)

	err = config.SaveConfig(cfg, "config_updated.json")
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	fmt.Println("Config updated successfully!")
}