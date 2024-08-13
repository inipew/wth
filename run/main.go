// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"run/cmd"
// 	"run/internal"

// 	"github.com/spf13/cobra"
// )

// var configFile = "./config.ini"
// var ctx context.Context
// var cancel context.CancelFunc

// func main() {
// 	ctx, cancel = context.WithCancel(context.Background())
// 	defer cancel()

// 	rootCmd := &cobra.Command{
// 		Use: "run",
// 	}

// 	// Register commands
// 	rootCmd.AddCommand(cmd.NewListCommand())
// 	rootCmd.AddCommand(cmd.NewExecuteCommand())
// 	rootCmd.AddCommand(cmd.NewAPICommand())

// 	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.ini", "Path to the INI configuration file")

// 	// Start file watcher in a separate goroutine
// 	go func() {
// 		err := internal.WatchConfig(ctx, configFile, func(cfg *internal.Config) {
// 			fmt.Println("Configuration file changed. Reloading...")
// 			// cfg, err := internal.LoadConfig(configFile)
// 			// if err != nil {
// 			// 	fmt.Printf("Error reloading config: %v\n", err)
// 			// 	return
// 			// }
// 			internal.DisplayCommands(cfg)
// 		})
// 		if err != nil {
// 			fmt.Printf("Error watching config: %v\n", err)
// 		}
// 	}()

//		// Execute the root command
//		if err := rootCmd.Execute(); err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//	}
package main

import (
	"context"
	"os"
	"run/cmd"
	"run/internal"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const numWorkers = 3

var configFile = "./config.ini"
var ctx context.Context
var cancel context.CancelFunc

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
		go internal.Worker(ctx, wg, taskQueue)
	}
}

func startConfigWatcher(ctx context.Context, configFile string, taskQueue chan func()) {
	if err := internal.WatchConfig(ctx, configFile, func(cfg *internal.Config) {
		taskQueue <- func() {
			logrus.Println("Configuration file changed. Reloading...")
			internal.DisplayCommands(cfg)
		}
	}); err != nil {
		logrus.Errorf("Error watching config: %v", err)
	}
}

