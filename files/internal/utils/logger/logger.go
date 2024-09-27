package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// LogLevel represents the logging level
type LogLevel string

const (
	// Debug level for verbose logging
	Debug LogLevel = "debug"
	// Info level for general operational entries
	Info LogLevel = "info"
	// Warn level for non-critical entries that deserve eyes
	Warn LogLevel = "warn"
	// Error level for errors that should definitely be noted
	Error LogLevel = "error"
	// Fatal level for very severe errors
	Fatal LogLevel = "fatal"
)

// Config holds the configuration for the logger
type Config struct {
	Level      LogLevel
	Output     io.Writer
	TimeFormat string
	NoColor    bool
	UseJSON    bool
}

// DefaultConfig returns a default configuration for the logger
func DefaultConfig() Config {
	return Config{
		Level:      Info,
		Output:     os.Stdout,
		TimeFormat: zerolog.TimeFormatUnix,
		NoColor:    false,
		UseJSON:    false,
	}
}

// Logger wraps zerolog.Logger to provide a custom interface
type Logger struct {
	zl zerolog.Logger
}

// New creates a new Logger instance with the given configuration
func New(cfg Config) *Logger {
	level, err := zerolog.ParseLevel(string(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = cfg.TimeFormat
	zerolog.SetGlobalLevel(level)

	var output io.Writer = cfg.Output
	if !cfg.UseJSON && !cfg.NoColor {
		output = zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: cfg.TimeFormat,
			NoColor:    cfg.NoColor,
		}
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return &Logger{zl: logger}
}

// Debug logs a message at the debug level
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log(Debug, msg, fields...)
}

// Info logs a message at the info level
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log(Info, msg, fields...)
}

// Warn logs a message at the warn level
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log(Warn, msg, fields...)
}

// Error logs a message at the error level
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log(Error, msg, fields...)
}

// Fatal logs a message at the fatal level and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log(Fatal, msg, fields...)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, msg string, fields ...interface{}) {
	zerologLevel, err := zerolog.ParseLevel(strings.ToLower(string(level)))
	if err != nil {
		zerologLevel = zerolog.InfoLevel
	}

	event := l.zl.WithLevel(zerologLevel)

	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fmt.Sprint(fields[i]), fields[i+1])
		}
	}

	event.Msg(msg)
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{zl: l.zl.With().Interface(key, value).Logger()}
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.zl.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{zl: ctx.Logger()}
}