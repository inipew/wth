package main

import (
	"fmt"
	"singconfig/internal/config"
	"singconfig/internal/config/inbound"
)

func main() {
	// config := config.BuildSingBoxConfig()
	// if err := config.SaveToFile("config.json"); err != nil {
	// 	log.Fatalf("Error saving JSON file: %s", err.Error())
	// }
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Example usage:
	// Add user
	newUser := inbound.UserConfig{Name: "pew", UUID: "new-uuid-here"}
	config.AddUser(cfg, newUser, "all", "")

	// Remove user
	config.RemoveUser(cfg, "default")

	// Update DNS
	config.UpdateDNS(cfg, "dns.google", "https", "prefer_ipv6")

	err = config.SaveConfig(cfg, "config_updated.json")
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	fmt.Println("Config updated successfully!")
}