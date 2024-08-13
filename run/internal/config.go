package internal

import (
	"fmt"
	"strings"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Paths    map[string]string
	Commands []Command
}

type Command struct {
	Name        string
	Value       string
	Command     string
	Description string
}

// LoadConfig loads the configuration from the specified INI file.
func LoadConfig(filename string) (*Config, error) {
	logrus.Infof("Memuat konfigurasi dari: %s", filename)
	cfg, err := ini.Load(filename)
	if err != nil {
		logrus.Errorf("Gagal memuat file INI: %v", err)
		return nil, fmt.Errorf("Gagal memuat file INI: %v", err)
	}

	config := &Config{
		Paths:    make(map[string]string),
		Commands: []Command{},
	}

	if err := parsePaths(cfg, config); err != nil {
		logrus.Errorf("Gagal memproses bagian [paths]: %v", err)
		return nil, fmt.Errorf("Gagal memproses bagian [paths]: %v", err)
	}

	if err := parseCommands(cfg, config); err != nil {
		logrus.Errorf("Gagal memproses bagian [commands]: %v", err)
		return nil, fmt.Errorf("Gagal memproses bagian [commands]: %v", err)
	}

	logrus.Info("Konfigurasi berhasil dimuat.")
	return config, nil
}

// parsePaths parses the [paths] section from the INI file and populates the Config.
func parsePaths(cfg *ini.File, config *Config) error {
	pathsSection := cfg.Section("paths")
	if pathsSection == nil {
		return fmt.Errorf("Bagian [paths] tidak ditemukan di file INI")
	}

	for _, key := range pathsSection.Keys() {
		if key.Name() == "" || key.String() == "" {
			return fmt.Errorf("Entri path tidak valid di bagian [paths]: kunci=%s, nilai=%s", key.Name(), key.String())
		}
		config.Paths[key.Name()] = key.String()
		logrus.Infof("Path dimuat: %s = %s", key.Name(), key.String())
	}

	return nil
}

// parseCommands parses the [commands] section from the INI file and populates the Commands slice.
func parseCommands(cfg *ini.File, config *Config) error {
	commandsSection := cfg.Section("commands")
	if commandsSection == nil {
		return fmt.Errorf("Bagian [commands] tidak ditemukan di file INI")
	}

	commandNames := make(map[string]struct{})

	for _, key := range commandsSection.Keys() {
		if strings.HasSuffix(key.Name(), "Name") {
			commandName := strings.TrimSuffix(key.Name(), "Name")
			if _, exists := commandNames[commandName]; !exists {
				command := Command{
					Name:        commandsSection.Key(commandName + "Name").String(),
					Value:       commandsSection.Key(commandName + "Value").String(),
					Command:     commandsSection.Key(commandName + "Command").String(),
					Description: commandsSection.Key(commandName + "Description").String(),
				}
				if command.Name == "" || command.Command == "" {
					return fmt.Errorf("Entri perintah tidak valid untuk: %s", commandName)
				}
				config.Commands = append(config.Commands, command)
				commandNames[commandName] = struct{}{}
				logrus.Infof("Perintah dimuat: %s = %s", command.Name, command.Command)
			}
		}
	}

	return nil
}

// ReplacePaths replaces placeholders in the command with corresponding values from paths map.
func ReplacePaths(command string, paths map[string]string) string {
	for key, value := range paths {
		placeholder := "%" + key + "%"
		command = strings.ReplaceAll(command, placeholder, value)
		logrus.Infof("Mengganti placeholder %s dengan %s", placeholder, value)
	}
	return command
}
