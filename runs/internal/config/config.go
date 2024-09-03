package config

import (
	"bytes"
	"fmt"
	"runs/internal/logger"
	"strings"
	"sync"
	"text/template"

	"gopkg.in/ini.v1"
)

// Config structure
type Config struct {
	WebConf   WebConf
	Paths     map[string]string
	Commands  []Command
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

var (
	ConfigData    Config
	ConfigMutex   sync.RWMutex
	ConfigCache   *Config
)

// Load configuration from file with mutex protection
func LoadConfig(configFilePath string) error {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	// If config is already loaded, skip reloading
	if ConfigCache != nil {
		return nil
	}

	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return err
	}

	ConfigData.WebConf = WebConf{
		Port:        cfg.Section("webconf").Key("port").String(),
		RestrictDir: cfg.Section("webconf").Key("restrictDir").MustBool(false),
	}

	// Initialize Paths map
	ConfigData.Paths = make(map[string]string)
	if err := parsePaths(cfg, &ConfigData); err != nil {
		return err
	}

	// Initialize Commands slice
	ConfigData.Commands = nil
	if err := parseCommands(cfg, &ConfigData); err != nil {
		return err
	}

	// Cache the loaded configuration
	ConfigCache = &ConfigData

	return nil
}

// parsePaths parses the [paths] section from the INI file and populates the Paths map.
func parsePaths(cfg *ini.File, config *Config) error {
	pathsSection := cfg.Section("paths")
	if pathsSection == nil {
		return fmt.Errorf("bagian [paths] tidak ditemukan di file INI")
	}

	for _, key := range pathsSection.Keys() {
		if key.Name() != "" && key.String() != "" {
			config.Paths[key.Name()] = key.String()
		} else {
			return fmt.Errorf("entri path tidak valid di bagian [paths]: kunci=%s, nilai=%s", key.Name(), key.String())
		}
	}

	return nil
}

// parseCommands parses the [commands] section from the INI file and populates the Commands slice.
func parseCommands(cfg *ini.File, config *Config) error {
	commandsSection := cfg.Section("commands")
	if commandsSection == nil {
		return fmt.Errorf("bagian [commands] tidak ditemukan di file INI")
	}

	var commandNames []string

	for _, key := range commandsSection.Keys() {
		if strings.HasSuffix(key.Name(), "Name") {
			commandName := strings.TrimSuffix(key.Name(), "Name")
			if !contains(commandNames, commandName) {
				command := Command{
					Name:        commandsSection.Key(commandName + "Name").String(),
					Value:       commandsSection.Key(commandName + "Value").String(),
					Command:     ReplacePlaceholders(commandsSection.Key(commandName + "Command").String(), config.Paths),
					Description: commandsSection.Key(commandName + "Description").String(),
				}
				if command.Name == "" || command.Command == "" {
					return fmt.Errorf("entri perintah tidak valid untuk: %s", commandName)
				}
				config.Commands = append(config.Commands, command)
				commandNames = append(commandNames, commandName)
			}
		}
	}

	return nil
}

// ReplacePlaceholders replaces placeholders in the command string with actual values from the paths map.
func ReplacePlaceholders(command string, paths map[string]string) string {
	tmpl, err := template.New("command").Parse(command)
	if err != nil {
		logger.Logger.Error().Str("command", command).Err(err).Msg("Error parsing template")
		return command
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, paths)
	if err != nil {
		logger.Logger.Error().Str("command", command).Err(err).Msg("Error executing template")
		return command
	}

	return buf.String()
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FindCommandByValue finds a command by its value.
func FindCommandByValue(value string) (*Command, error) {
	ConfigMutex.RLock()
	defer ConfigMutex.RUnlock()

	for _, cmd := range ConfigData.Commands {
		if cmd.Value == value {
			return &cmd, nil
		}
	}

	return nil, fmt.Errorf("command with value %s not found", value)
}

// PrintConfig prints the current configuration.
func PrintConfig() {
	ConfigMutex.RLock()
	defer ConfigMutex.RUnlock()

	fmt.Println("\n=== Configuration Loaded ===")
	fmt.Printf("WebConf:\n")
	fmt.Printf("  Port: %s\n", ConfigData.WebConf.Port)
	fmt.Printf("  RestrictDir: %v\n", ConfigData.WebConf.RestrictDir)

	fmt.Printf("\nPaths:\n")
	for key, value := range ConfigData.Paths {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Printf("\nCommands:\n")
	for _, cmd := range ConfigData.Commands {
		fmt.Printf("  %s:\n", cmd.Value)
		fmt.Printf("    Name: %s\n", cmd.Name)
		fmt.Printf("    Value: %s\n", cmd.Value)
		fmt.Printf("    Command: %s\n", cmd.Command)
		fmt.Printf("    Description: %s\n", cmd.Description)
	}
}
