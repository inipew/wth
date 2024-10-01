package main

import (
	"fmt"
	"log"
	"singconfig/internal/account"
	"singconfig/internal/config"
	"singconfig/internal/utils"
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
	uuid := utils.GenerateUUID()
	account.AddUser(cfg, "httpupgrade", "", "akupew", uuid)
	account.AddUser(cfg, "ws", "", "akupew", uuid)
	config.ModifyLogLevel(cfg,"debug")


	// Remove user
	// config.RemoveUser(cfg, "default")

	// Update DNS
	// config.UpdateDNS(cfg, "dns.google", "https", "prefer_ipv6")
	config.DisplayInboundDetails(cfg)
	account.PrintInboundConfigs(cfg)

	err = config.WriteConfigFile(cfg, "config_updated.json")
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}

	fmt.Println("Config updated successfully!")
}