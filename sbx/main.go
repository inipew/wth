package main

import (
	"fmt"
	"sbx/internal/cmd"
	"sbx/internal/commands"
	"sbx/internal/config"
	"sbx/internal/logger"
	"sbx/internal/utils"

	"github.com/rs/zerolog"
)

func main() {
	cfg := logger.DefaultConfig()
    cfg.LogLevel = zerolog.DebugLevel
    // cfg.OutputFile = "app.log"
	cfg.UseColor = true

    err := logger.InitGlobalLogger(cfg)
    if err != nil {
        panic(err)
    }
    log := logger.GetLogger()

	if !utils.CheckRoot() {
        log.Fatal().Msg("This program must be run as root")
    }


	config.CreateAllDirs()

	myCLI := cmd.NewCLI("1.0.0", "SBX - A toolbox for server management")

	// Register commands
	myCLI.AddCommand(commands.CreateSingCommand())
	myCLI.AddCommand(commands.CreateCaddyCommand())
	myCLI.AddCommand(commands.CreateCaddyfileCommand())
	myCLI.AddCommand(commands.CreateInstallCommand())
	myCLI.AddCommand(commands.CreateServiceCommand())

	myCLI.SetRootCommand(createRootCommand())

	// Execute the CLI
	myCLI.Execute()
}

func createRootCommand() *cmd.Command {
	return &cmd.Command{
		Name:        "sbx",
		Description: "SBX root command - displays general information",
		Run: func(cmd *cmd.Command, args []string) error {
			fmt.Println("Welcome to SBX - Your Server Management Toolbox")
			fmt.Println("Use 'sbx help' to see available commands.")
			return nil
		},
	}
}