package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func DisplayCommands(cfg *Config) {
	fmt.Println("Available Commands:")
	for _, cmd := range cfg.Commands {
		fmt.Printf("%s - %s\n", cmd.Name, cmd.Description)
	}
}

// RunCommand executes a command with a timeout and returns the output and any error
func RunCommand(ctx context.Context, command string, timeout time.Duration) (string, error) {
	// Set a timeout for the command
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Determine the command to run based on the operating system
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd.exe", "/C", command)
	case "linux", "darwin", "android":
		sh := "sh"
		if os == "android" {
			sh = "su"
		}
		cmd = exec.CommandContext(ctx, sh, "-c", command)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", os)
	}

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	
	// Check for timeout error
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command timed out: %s", command)
	}

	// Handle the different error cases more specifically
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("command not found: %s", command)
		}
		if errors.Is(err, os.ErrPermission) {
			return "", fmt.Errorf("permission denied: %s", command)
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("command exited with code %d: %v", exitErr.ExitCode(), exitErr)
		}
		return "", fmt.Errorf("error running command '%s': %v", command, err)
	}

	return string(output), nil
}

func RunCommandDirectly(command string, args []string) error {
    var execErr error
    switch runtime.GOOS {
    case "windows":
        execErr = syscall.Exec("cmd.exe", append([]string{"/C", command}, args...), os.Environ())
    case "linux", "darwin", "android":
        execErr = syscall.Exec("/bin/sh", append([]string{"-c", command}, args...), os.Environ())
    default:
        return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }
    return execErr
}

// Worker is a worker function that processes tasks from the taskQueue.
// func Worker(ctx context.Context, wg *sync.WaitGroup, taskQueue chan func()) {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case task, ok := <-taskQueue:
// 			if !ok {
// 				return // taskQueue is closed
// 			}
// 			task()
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
func Worker(ctx context.Context, wg *sync.WaitGroup, taskQueue chan func()) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return // Context is canceled or timed out
		case task, ok := <-taskQueue:
			if !ok {
				return // taskQueue is closed
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Handle or log panic
						logrus.Error("Error panic")
					}
				}()
				task() // Execute the task
			}()
		}
	}
}
