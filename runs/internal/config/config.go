package config

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"

	"github.com/rs/zerolog"
	"gopkg.in/ini.v1"
)

// ConfigProvider defines the interface for loading configuration
type ConfigProvider interface {
	Load(ctx context.Context) (*Config, error)
}

// Config structure
type Config struct {
	WebConf  WebConf
	Paths    map[string]string
	Commands []Command
}

// WebConf structure
type WebConf struct {
	Port        string
	RestrictDir bool
}

// Command structure
type Command struct {
	Name        string
	Value       string
	Command     string
	Description string
}

// FileConfigProvider implements ConfigProvider for INI files
type FileConfigProvider struct {
	FilePath string
	logger   zerolog.Logger
}

// NewFileConfigProvider creates a new FileConfigProvider
func NewFileConfigProvider(filePath string, logger zerolog.Logger) *FileConfigProvider {
	return &FileConfigProvider{
		FilePath: filePath,
		logger:   logger,
	}
}

// Load implements ConfigProvider.Load for INI files
func (fcp *FileConfigProvider) Load(ctx context.Context) (*Config, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, fcp.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	config := &Config{
		WebConf: WebConf{
			Port:        cfg.Section("webconf").Key("port").MustString(":8080"),
			RestrictDir: cfg.Section("webconf").Key("restrictDir").MustBool(false),
		},
		Paths:    make(map[string]string),
		Commands: []Command{},
	}

	if err := parsePaths(cfg, config); err != nil {
		return nil, fmt.Errorf("failed to parse paths: %w", err)
	}

	if err := parseCommands(cfg, config); err != nil {
		return nil, fmt.Errorf("failed to parse commands: %w", err)
	}

	return config, nil
}

func parsePaths(cfg *ini.File, config *Config) error {
	pathsSection := cfg.Section("paths")
	for _, key := range pathsSection.Keys() {
		if key.Name() != "" && key.String() != "" {
			config.Paths[key.Name()] = key.String()
		} else {
			return fmt.Errorf("invalid path entry in [paths] section: key=%s, value=%s", key.Name(), key.String())
		}
	}
	return nil
}

func parseCommands(cfg *ini.File, config *Config) error {
	commandsSection := cfg.Section("commands")
	for _, key := range commandsSection.Keys() {
		if key.Name() == "" || len(key.Name()) < 4 || key.Name()[len(key.Name())-4:] != "Name" {
			continue
		}
		
		baseName := key.Name()[:len(key.Name())-4]
		command := Command{
			Name:        key.String(),
			Value:       commandsSection.Key(baseName + "Value").String(),
			Command:     replacePlaceholders(commandsSection.Key(baseName + "Command").String(), config.Paths),
			Description: commandsSection.Key(baseName + "Description").String(),
		}
		
		if command.Name == "" || command.Value == "" || command.Command == "" {
			return fmt.Errorf("invalid command entry for: %s", baseName)
		}
		
		config.Commands = append(config.Commands, command)
	}
	return nil
}

func replacePlaceholders(command string, paths map[string]string) string {
	tmpl, err := template.New("command").Parse(command)
	if err != nil {
		return command
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, paths)
	if err != nil {
		return command
	}

	return buf.String()
}

// ConfigManager manages the configuration
type ConfigManager struct {
	provider ConfigProvider
	config   *Config
	mu       sync.RWMutex
	logger   zerolog.Logger
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(provider ConfigProvider, logger zerolog.Logger) *ConfigManager {
	return &ConfigManager{
		provider: provider,
		logger:   logger,
	}
}

// Load loads the configuration
func (cm *ConfigManager) Load(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, err := cm.provider.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cm.config = config
	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// FindCommandByValue finds a command by its value
func (cm *ConfigManager) FindCommandByValue(value string) (*Command, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, cmd := range cm.config.Commands {
		if cmd.Value == value {
			return &cmd, nil
		}
	}

	return nil, fmt.Errorf("command with value %s not found", value)
}

// LogConfig logs the current configuration
func (cm *ConfigManager) LogConfig() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cm.logger.Info().
		Str("port", cm.config.WebConf.Port).
		Bool("restrictDir", cm.config.WebConf.RestrictDir).
		Interface("paths", cm.config.Paths).
		Interface("commands", cm.config.Commands).
		Msg("Configuration loaded")
}