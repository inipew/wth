package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/fileutils"
	"sbx/internal/logger"
	"sync"
)

// AppConfig holds all configuration constants
type AppConfig struct {
	WorkDir              string
	BackupDir            string
	TmpDir               string
	LogDir               string
	DomainFilePath       string
	SingboxDir           string
	SingboxConfDir       string
	SingboxLogFilePath   string
	CaddyWorkDir         string
	CaddyFilePath        string
	CaddyLogFile         string
	CaddyAccessLogFile   string
	BinaryBinPath		 string
	SingboxBinPath       string
	CaddyBinPath         string
	CaddyServicePath     string
	SingboxServicePath   string
}

var (
	appConfig *AppConfig
	once      sync.Once
)

// GetConfig returns the singleton instance of AppConfig
func GetConfig() *AppConfig {
	once.Do(func() {
		appConfig = &AppConfig{
			WorkDir:              "/etc/wth",
			BackupDir:            "/etc/wth/backup",
			TmpDir:               "/etc/wth/tmp",
			LogDir:               "/etc/wth/log",
			DomainFilePath:       "/etc/wth/domain.txt",
			SingboxDir:           "/etc/wth/sing-box",
			SingboxConfDir:       "/etc/wth/sing-box/config",
			SingboxLogFilePath:   "/etc/wth/log/sing-box.log",
			CaddyWorkDir:         "/etc/wth/caddy",
			CaddyFilePath:        "/etc/wth/caddy/Caddyfile",
			CaddyLogFile:         "/etc/wth/log/caddy.log",
			CaddyAccessLogFile:   "/etc/wth/log/access.log",
			BinaryBinPath:		  "/usr/local/bin",
			SingboxBinPath:       "/usr/local/bin/sing-box",
			CaddyBinPath:         "/usr/local/bin/caddy",
			CaddyServicePath:     "/etc/systemd/system/caddy.service",
			SingboxServicePath:   "/etc/systemd/system/sing-box.service",
		}
	})
	return appConfig
}

// Helper functions for easy access to config values
func WorkDir() string              { return GetConfig().WorkDir }
func BackupDir() string            { return GetConfig().BackupDir }
func TmpDir() string               { return GetConfig().TmpDir }
func LogDir() string               { return GetConfig().LogDir }
func DomainFilePath() string       { return GetConfig().DomainFilePath }
func SingboxDir() string           { return GetConfig().SingboxDir }
func SingboxConfDir() string       { return GetConfig().SingboxConfDir }
func SingboxLogFilePath() string   { return GetConfig().SingboxLogFilePath }
func CaddyWorkDir() string         { return GetConfig().CaddyWorkDir }
func CaddyFilePath() string        { return GetConfig().CaddyFilePath }
func CaddyLogFile() string         { return GetConfig().CaddyLogFile }
func CaddyAccessLogFile() string   { return GetConfig().CaddyAccessLogFile }
func BinaryBinPath() string       { return GetConfig().BinaryBinPath }
func SingboxBinPath() string       { return GetConfig().SingboxBinPath }
func CaddyBinPath() string         { return GetConfig().CaddyBinPath }
func CaddyServicePath() string     { return GetConfig().CaddyServicePath }
func SingboxServicePath() string   { return GetConfig().SingboxServicePath }

// CreateAllDirs creates all necessary directories for configuration
func CreateAllDirs() error {
	cfg := GetConfig()
	dirs := []string{
		cfg.WorkDir,
		cfg.BackupDir,
		cfg.TmpDir,
		cfg.LogDir,
		cfg.SingboxDir,
		cfg.SingboxConfDir,
		cfg.CaddyWorkDir,
	}

	for _, dir := range dirs {
		if err := createDirIfNotExists(dir); err != nil {
			return err
		}
	}

	return createLogFiles()
}

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(dir string) error {
	exists, err := fileutils.IsDirectory(dir)
	if err != nil {
		return fmt.Errorf("error checking directory %s: %w", dir, err)
	}

	if !exists {
		if err := fileutils.CreateDir(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logger.GetLogger().Info().Msgf("Directory created: %s", dir)
	} else {
		logger.GetLogger().Debug().Msgf("Directory already exists: %s", dir)
	}

	return nil
}

// createLogFiles creates log files if they do not exist
func createLogFiles() error {
	cfg := GetConfig()
	logFiles := []string{
		cfg.SingboxLogFilePath,
		cfg.CaddyLogFile,
		cfg.CaddyAccessLogFile,
	}

	for _, logFile := range logFiles {
		if err := createLogFileIfNotExists(logFile); err != nil {
			return err
		}
	}

	return nil
}

// createLogFileIfNotExists creates a log file if it doesn't exist
func createLogFileIfNotExists(filePath string) error {
	exists, err := fileutils.IsFile(filePath)
	if err != nil {
		return fmt.Errorf("error checking log file %s: %w", filePath, err)
	}

	if !exists {
		if err := fileutils.CreateDir(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for log file %s: %w", filePath, err)
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to create log file %s: %w", filePath, err)
		}
		file.Close()

		logger.GetLogger().Info().Msgf("Log file created: %s", filePath)
	} else {
		logger.GetLogger().Debug().Msgf("Log file already exists: %s", filePath)
	}

	return nil
}