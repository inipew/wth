package main

import (
	"context"
	"os"
	"run/cmd"
	"run/internal/config"
	"run/internal/utils"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const numWorkers = 3

var configFile = "./config.ini"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rootCmd := setupCommands()

	// Initialize worker pool
	taskQueue := make(chan func(), 10)
	var wg sync.WaitGroup

	startWorkers(ctx, &wg, taskQueue, numWorkers)

	// Start file watcher in a separate goroutine
	go startConfigWatcher(ctx, configFile, taskQueue)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("Error executing command: %v", err)
		os.Exit(1)
	}

	close(taskQueue)
	wg.Wait() // Wait for all workers to finish
}
func setupCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "run",
	}

	// Register commands
	rootCmd.AddCommand(cmd.NewListCommand())
	rootCmd.AddCommand(cmd.NewExecuteCommand())
	rootCmd.AddCommand(cmd.NewAPICommand())

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.ini", "Path to the INI configuration file")
	return rootCmd
}

func startWorkers(ctx context.Context, wg *sync.WaitGroup, taskQueue chan func(), numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go utils.Worker(ctx, wg, taskQueue)
	}
}

func startConfigWatcher(ctx context.Context, configFile string, taskQueue chan func()) {
	if err := utils.WatchConfig(ctx, configFile, func(cfg *config.Config) {
		taskQueue <- func() {
			logrus.Println("Configuration file changed. Reloading...")
			utils.DisplayCommands(cfg)
		}
	}); err != nil {
		logrus.Errorf("Error watching config: %v", err)
	}
}

