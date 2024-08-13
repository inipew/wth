package main

import (
	"fmt"
	"os"
	"github.com/inipew/wth/caddy"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "caddy":
		if len(os.Args) < 3 {
			printCaddyUsage()
			return
		}

		subcommand := os.Args[2]
		switch subcommand {
		case "install":
			caddy.InstallCaddy()
		case "start":
			caddy.StartCaddy()
		case "restart":
			caddy.RestartCaddy()
		case "stop":
			caddy.StopCaddy()
		case "changedomain":
			if len(os.Args) < 4 {
				fmt.Println("Usage: ./binary caddy changedomain <newdomain>")
				return
			}
			newDomain := os.Args[3]
			caddy.ChangeDomain(newDomain)
		default:
			fmt.Println("Unknown subcommand:", subcommand)
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}

func printUsage() {
	fmt.Println("Usage: ./binary <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  caddy <subcommand> - Manage Caddy")
}

func printCaddyUsage() {
	fmt.Println("Usage: ./binary caddy <subcommand> [options]")
	fmt.Println("Subcommands:")
	fmt.Println("  install       - Install Caddy")
	fmt.Println("  start         - Start Caddy")
	fmt.Println("  restart       - Restart Caddy")
	fmt.Println("  stop          - Stop Caddy")
	fmt.Println("  changedomain  - Change domain in Caddyfile")
}
