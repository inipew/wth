package cmd

import (
	"context"
	"fmt"
	"run/internal"
	"time"

	"github.com/spf13/cobra"
)

var (
	configFile string
	timeout    = 10 * time.Second
)

// NewExecuteCommand creates a new execute command to run a command by its value
func NewExecuteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [value]",
		Short: "Execute a command from the configuration based on the Value",
		Run:   executeCommand,
	}

	// Define the config flag for the execute command
	// cmd.Flags().StringVarP(&configFile, "config", "c", "config.ini", "Path to the INI configuration file")
	return cmd
}

// executeCommand finds and runs a command based on the provided value
func executeCommand(cmd *cobra.Command, args []string) {
	cfg, err := internal.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	commandValue := args[0]
	command, found := findCommandByValue(cfg, commandValue)
	if !found {
		fmt.Printf("Command with Value '%s' not found\n", commandValue)
		return
	}

	// Replace paths placeholders in the command string
	command.Command = internal.ReplacePaths(command.Command, cfg.Paths)
	
	output, err := internal.RunCommand(context.Background(), command.Command, timeout)
	if err != nil {
		fmt.Printf("Failed to execute command: %v\n", err)
		return
	}

	fmt.Printf("Command Output:\n%s\n", output)
}

// findCommandByValue searches for a command by its value
func findCommandByValue(cfg *internal.Config, value string) (*internal.Command, bool) {
	for _, cmd := range cfg.Commands {
		if cmd.Value == value {
			return &cmd, true
		}
	}
	return nil, false
}
