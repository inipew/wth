package systemdservice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	cfg "sbx/internal/config"
	"sbx/internal/logger"
	"strings"
	"text/template"
)

// ServiceConfig stores the configuration for a systemd service.
type ServiceConfig struct {
	Description        string
	Documentation      string
	After              []string
	Requires           []string
	ExecStart          string
	ExecReload         string
	TimeoutStopSec     string
	Restart            string
	RestartSec         string
	LimitNOFILE        string
	AmbientCaps        []string
	CapabilityBounding []string
	PrivateTmp         bool
	ProtectSystem      string
}

const serviceTemplate = `[Unit]
Description={{ .Description }}
Documentation={{ .Documentation }}
After={{ join .After " " }}
{{ if .Requires }}Requires={{ join .Requires " " }}{{ end }}

[Service]
ExecStart={{ .ExecStart }}
ExecReload={{ .ExecReload }}
LimitNOFILE={{ .LimitNOFILE }}
AmbientCapabilities={{ join .AmbientCaps " " }}
{{ if .CapabilityBounding }}CapabilityBoundingSet={{ join .CapabilityBounding " " }}{{ end }}
{{ if .TimeoutStopSec }}TimeoutStopSec={{ .TimeoutStopSec }}{{ end }}
{{ if .Restart }}Restart={{ .Restart }}{{ end }}
{{ if .RestartSec }}RestartSec={{ .RestartSec }}{{ end }}
{{ if .PrivateTmp }}PrivateTmp=true{{ end }}
{{ if .ProtectSystem }}ProtectSystem={{ .ProtectSystem }}{{ end }}

[Install]
WantedBy=multi-user.target
`

// ServiceManager handles the creation and management of systemd services
type ServiceManager struct {
	tmpl *template.Template
}

// NewServiceManager creates a new ServiceManager
func NewServiceManager() (*ServiceManager, error) {
	tmpl, err := template.New("service").Funcs(template.FuncMap{"join": strings.Join}).Parse(serviceTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service template: %w", err)
	}
	return &ServiceManager{tmpl: tmpl}, nil
}

// CreateServiceFile creates a systemd service file based on the provided configuration.
func (sm *ServiceManager) CreateServiceFile(serviceName string, config ServiceConfig) error {
	filePath := filepath.Join("/etc/systemd/system", fmt.Sprintf("%s.service", serviceName))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create service file at %s: %w", filePath, err)
	}
	defer file.Close()

	if err := sm.tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	logger.GetLogger().Info().Str("Service", serviceName).Msg("Service file generated successfully")
	return nil
}

// GenerateCaddyService generates the caddy.service file.
func (sm *ServiceManager) GenerateCaddyService() error {
	config := ServiceConfig{
		Description:    "Caddy",
		Documentation:  "https://caddyserver.com/docs/",
		After:          []string{"network.target", "network-online.target"},
		Requires:       []string{"network-online.target"},
		ExecStart:      fmt.Sprintf("%s run --environ --config %s", cfg.CaddyBinPath(), cfg.CaddyFilePath()),
		ExecReload:     fmt.Sprintf("%s reload --config %s --force", cfg.CaddyBinPath(), cfg.CaddyFilePath()),
		TimeoutStopSec: "5s",
		LimitNOFILE:    "1048576",
		PrivateTmp:     true,
		ProtectSystem:  "full",
		AmbientCaps:    []string{"CAP_NET_ADMIN", "CAP_NET_BIND_SERVICE"},
	}
	return sm.CreateServiceFile("caddy", config)
}

// GenerateSingBoxService generates the sing-box.service file.
func (sm *ServiceManager) GenerateSingBoxService() error {
	config := ServiceConfig{
		Description:        "sing-box service",
		Documentation:      "https://sing-box.sagernet.org",
		After:              []string{"network.target", "nss-lookup.target", "network-online.target"},
		ExecStart:          fmt.Sprintf("%s -D /var/lib/sing-box -C %s run", cfg.SingboxBinPath(), cfg.SingboxConfDir()),
		ExecReload:         "/bin/kill -HUP $MAINPID",
		Restart:            "on-failure",
		RestartSec:         "10s",
		LimitNOFILE:        "infinity",
		CapabilityBounding: []string{"CAP_NET_ADMIN", "CAP_NET_BIND_SERVICE", "CAP_SYS_PTRACE", "CAP_DAC_READ_SEARCH"},
		AmbientCaps:        []string{"CAP_NET_ADMIN", "CAP_NET_BIND_SERVICE", "CAP_SYS_PTRACE", "CAP_DAC_READ_SEARCH"},
	}
	return sm.CreateServiceFile("sing-box", config)
}

// ServiceOperation defines the type for systemd service operations
type ServiceOperation func(string) error

// systemctlCommand executes a systemctl command and returns an error if any occurs.
func systemctlCommand(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s failed: %w, output: %s", strings.Join(args, " "), err, output)
	}
	return nil
}

// EnableService enables a systemd service.
func EnableService(serviceName string) error {
	return systemctlCommand("enable", serviceName)
}

// DisableService disables a systemd service.
func DisableService(serviceName string) error {
	return systemctlCommand("disable", serviceName)
}

// StartService starts a systemd service.
func StartService(serviceName string) error {
	return systemctlCommand("start", serviceName)
}

// StopService stops a systemd service.
func StopService(serviceName string) error {
	return systemctlCommand("stop", serviceName)
}

// RestartService restarts a systemd service.
func RestartService(serviceName string) error {
	return systemctlCommand("restart", serviceName)
}