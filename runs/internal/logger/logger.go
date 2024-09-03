package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// Logger instance
var Logger zerolog.Logger

// InitLogger initializes the logger
func InitLogger() error {
	// Create a new logger instance
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	
	// Check if any error occurs during logger setup
	if logger.GetLevel() == zerolog.NoLevel {
		return fmt.Errorf("failed to initialize logger: invalid log level")
	}
	
	Logger = logger
	return nil
}
