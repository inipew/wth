package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger with console writer
func InitLogger() error {
	// Set up the console writer with desired settings
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf(" | [ %s ]", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf(" %s", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%v", i)
		},
	}

	// Create the logger with the console writer
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	if logger.GetLevel() == zerolog.NoLevel {
		return fmt.Errorf("failed to initialize logger: invalid log level")
	}

	// Set global logger
	log.Logger = logger

	return nil
}
