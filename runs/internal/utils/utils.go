package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"runs/internal/config"
	"runs/internal/logger"
)

// ErrUnsupportedOS is returned when the operating system is not supported
var ErrUnsupportedOS = errors.New("unsupported operating system")

// DisplayCommands prints available commands to stdout
func DisplayCommands(cfg *config.Config) {
	if cfg == nil {
		logger.GetLogger().Info().Msg("Error: Configuration is nil")
		return
	}

	fmt.Println("Available Commands:")
	for _, cmd := range cfg.Commands {
		fmt.Printf("%-15s - %s\n", cmd.Name, cmd.Description)
	}
}

// RunCommand executes a command with a timeout and returns the output and any error
func RunCommand(ctx context.Context, command string, timeout time.Duration) (string, error) {
	if ctx == nil {
		return "", errors.New("context is nil")
	}

	if command == "" {
		return "", errors.New("command is empty")
	}

	if timeout <= 0 {
		return "", errors.New("timeout must be positive")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd, err := createOSSpecificCommand(ctx, command)
	if err != nil {
		return "", fmt.Errorf("failed to create command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", handleCommandError(ctx, command, err, output)
	}

	return strings.TrimSpace(string(output)), nil
}

func createOSSpecificCommand(ctx context.Context, command string) (*exec.Cmd, error) {
	switch os := runtime.GOOS; os {
	case "windows":
		return exec.CommandContext(ctx, "cmd.exe", "/C", command), nil
	case "linux", "darwin":
		return exec.CommandContext(ctx, "sh", "-c", command), nil
	case "android":
		return exec.CommandContext(ctx, "su", "-c", command), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOS, os)
	}
}

func handleCommandError(ctx context.Context, command string, err error, output []byte) error {
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out: %s", command)
	}

	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("command not found: %s", command)
	}

	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("permission denied: %s", command)
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("command exited with code %d: %s\nOutput: %s", exitErr.ExitCode(), command, string(output))
	}

	return fmt.Errorf("error running command '%s': %w\nOutput: %s", command, err, string(output))
}