package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Files  FilesConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port              string
	UseEmbeddedFiles  bool
	MaxUploadSize     int64
	TempUploadDir     string
	ReadTimeout       int
	WriteTimeout      int
	GracefulShutdown  int
}

// FilesConfig holds file-related configuration
type FilesConfig struct {
	StorageDir     	string
	ShowHiddenFiles	bool
	MaxFileSize    	int64
	ArchiveEnabled 	bool
}

// Load reads the configuration file and returns a Config struct
func Load(filename string) (*Config, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	config := &Config{}
	err = cfg.MapTo(config)
	if err != nil {
		return nil, fmt.Errorf("failed to map config: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	config.setDefaults()

	return config, nil
}

func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Files.StorageDir == "" {
		return fmt.Errorf("storage directory is required")
	}

	if _, err := os.Stat(c.Files.StorageDir); os.IsNotExist(err) {
		return fmt.Errorf("storage directory does not exist: %s", c.Files.StorageDir)
	}

	return nil
}

func (c *Config) setDefaults() {
	if c.Server.Port == "" {
		c.Server.Port = ":5678"
	}

	if c.Server.MaxUploadSize == 0 {
		c.Server.MaxUploadSize = 100 * 1024 * 1024 // 100 MB
	}

	if c.Server.TempUploadDir == "" {
		c.Server.TempUploadDir = os.TempDir()
	}

	if c.Server.ReadTimeout == 0 {
		c.Server.ReadTimeout = 30 // 30 seconds
	}

	if c.Server.WriteTimeout == 0 {
		c.Server.WriteTimeout = 30 // 30 seconds
	}

	if c.Server.GracefulShutdown == 0 {
		c.Server.GracefulShutdown = 15 // 15 seconds
	}

	if c.Files.MaxFileSize == 0 {
		c.Files.MaxFileSize = 100 * 1024 * 1024 // 50 MB
	}
}

// GetAbsoluteStoragePath returns the absolute path of the storage directory
func (c *Config) GetAbsoluteStoragePath() (string, error) {
	return filepath.Abs(c.Files.StorageDir)
}