package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	globalLogger *zerolog.Logger
	once         sync.Once
)

// Config holds the configuration for the logger
type Config struct {
	UseColor    bool
	LogLevel    zerolog.Level
	TimeFormat  string
	OutputFile  string // If empty, logs to stdout
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		UseColor:    true,
		LogLevel:    zerolog.InfoLevel,
		TimeFormat:  zerolog.TimeFormatUnix,
		OutputFile:  "", // Default to stdout
	}
}

// InitGlobalLogger initializes the global logger with the given configuration
func InitGlobalLogger(cfg Config) error {
	var err error
	once.Do(func() {
		var output zerolog.ConsoleWriter
		if cfg.OutputFile != "" {
			file, openErr := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if openErr != nil {
				err = fmt.Errorf("failed to open log file: %w", openErr)
				return
			}
			output = zerolog.ConsoleWriter{Out: file, NoColor: true}
		} else {
			output = zerolog.ConsoleWriter{
				Out:     os.Stdout,
				NoColor: !cfg.UseColor,
			}
		}

		output.TimeFormat = cfg.TimeFormat
		output.FormatLevel = formatLevel
		output.FormatMessage = formatMessage
		output.FormatFieldName = formatFieldName
		output.FormatFieldValue = formatFieldValue

		logger := zerolog.New(output).Level(cfg.LogLevel).With().Timestamp().Logger()
		globalLogger = &logger
	})

	return err
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	if globalLogger == nil {
		// If not initialized, create a default logger
		defaultLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		return &defaultLogger
	}
	return globalLogger
}

// Helper functions for formatting
func formatLevel(i interface{}) string {
	return fmt.Sprintf("| %-6s|", i)
}

func formatMessage(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func formatFieldName(i interface{}) string {
	return fmt.Sprintf("%s:", i)
}

func formatFieldValue(i interface{}) string {
	return fmt.Sprintf("%v", i)
}