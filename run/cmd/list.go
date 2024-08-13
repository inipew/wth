package cmd

import (
	"fmt"
	"run/internal"

	"github.com/spf13/cobra"
)

// NewListCommand creates a new list command to display available commands
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the command list",
		Run:   listCommands,
	}

	// Define the config flag for the list command
	cmd.Flags().StringVarP(&configFile, "config", "c", "config.ini", "Path to the INI configuration file")
	return cmd
}

// listCommands loads the configuration and displays the list of commands
func listCommands(cmd *cobra.Command, args []string) {
	cfg, err := internal.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	internal.DisplayCommands(cfg)
}
