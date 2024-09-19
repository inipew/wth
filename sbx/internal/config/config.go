package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/fileutils"
)

// Constants for directories and paths
const (
	WorkDir              = "/etc/wth"
	BackupDir            = WorkDir + "/backup"
	TmpDir               = WorkDir + "/tmp"
	LogDir               = WorkDir + "/log"
	DomainFilePath       = WorkDir + "/domain.txt"
	SingboxDir           = WorkDir + "/sing-box"
	SingboxConfDir       = WorkDir + "/sing-box/config"
	SingboxLogFilePath   = LogDir + "/sing-box.log"
	CaddyWorkDir         = WorkDir + "/caddy"
	CaddyFilePath        = CaddyWorkDir + "/Caddyfile"
	CaddyLogFile         = LogDir + "/caddy.log"
	CaddyAccessLogFile   = LogDir + "/access.log"

	SingboxBinPath       = "/usr/local/bin/sing-box"
	CaddyBinPath         = "/usr/local/bin/caddy"
	CaddyServicePath     = "/etc/systemd/system/caddy.service"
	SingboxServicePath    = "/etc/systemd/system/sing-box.service"
)

// CreateAllDirs creates all necessary directories for configuration
func CreateAllDirs() error {
	dirs := []string{
		WorkDir,
		BackupDir,
		TmpDir,
		LogDir,
		SingboxDir,
		SingboxConfDir,
		CaddyWorkDir,
	}

	for _, dir := range dirs {
		if err := fileutils.CreateDir(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	if err := createLogFiles(); err != nil {
		return err
	}

	return nil
}

// createLogFiles creates log files if they do not exist
func createLogFiles() error {
	logFiles := []string{
		SingboxLogFilePath,
		CaddyLogFile,
		CaddyAccessLogFile,
	}

	for _, logFile := range logFiles {
		if err := createLogFile(logFile); err != nil {
			return fmt.Errorf("failed to create log file %s: %w", logFile, err)
		}
	}

	return nil
}

// createLogFile creates a log file if it does not exist
func createLogFile(filePath string) error {
	if fileutils.Exists(filePath) {
		return nil // File already exists
	}

	// Create the directory for the log file if it does not exist
	dir := filepath.Dir(filePath)
	if err := fileutils.CreateDir(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for log file %s: %w", dir, err)
	}

	// Create an empty log file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file %s: %w", filePath, err)
	}
	defer file.Close()

	return nil
}
