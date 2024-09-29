package caddyfile

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sbx/internal/config"
	"sbx/internal/logger"
	"text/template"
)

// Config holds the information needed to generate the Caddyfile
type Config struct {
	Domain string
}

// CaddyfileGenerator handles the generation of Caddyfiles
type CaddyfileGenerator struct {
	templateContent string
	config          Config
}

// NewCaddyfileGenerator creates a new CaddyfileGenerator
func NewCaddyfileGenerator(templateContent string, config Config) *CaddyfileGenerator {
	return &CaddyfileGenerator{
		templateContent: templateContent,
		config:          config,
	}
}

// Generate generates a Caddyfile based on the provided config
func (cg *CaddyfileGenerator) Generate() (string, error) {
	if err := cg.validateConfig(); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}

	tmpl, err := template.New("caddyfile").Parse(cg.templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, cg.prepareTemplateData()); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

// validateConfig checks if the domain is valid
func (cg *CaddyfileGenerator) validateConfig() error {
	if cg.config.Domain == "" {
		return errors.New("domain cannot be empty")
	}
	return nil
}

// prepareTemplateData prepares the data for the template
func (cg *CaddyfileGenerator) prepareTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Domain":          cg.config.Domain,
		"CaddyLogFile":    config.CaddyLogFile(),
		"TrustedProxies":  trustedProxies(),
		"WebSocketMap":    generateWebSocketMap(),
		"HttpUpgradeMap":  generateHttpUpgradeMap(),
	}
}

// trustedProxies returns a string representation of trusted proxies
func trustedProxies() string {
	return "173.245.48.0/20 103.21.244.0/22 103.22.200.0/22 103.31.4.0/22 " +
		"141.101.64.0/18 108.162.192.0/18 190.93.240.0/20 188.114.96.0/20 " +
		"197.234.240.0/22 198.41.128.0/17 162.158.0.0/15 104.16.0.0/13 " +
		"104.24.0.0/14 172.64.0.0/13 131.0.72.0/22 " +
		"2400:cb00::/32 2606:4700::/32 2803:f800::/32 2405:b500::/32 " +
		"2405:8100::/32 2a06:98c0::/29 2c0f:f248::/32"
}

// generateWebSocketMap generates the mapping for WebSocket
func generateWebSocketMap() string {
	return `map {path} {backend} {
		/trojan 127.0.0.1:8003
		/vmess 127.0.0.1:8002
		/vless 127.0.0.1:8001
		/trojan/Tun 127.0.0.1:8007
		/vmess/Tun 127.0.0.1:8008
		/vless/Tun 127.0.0.1:8009
	}`
}

// generateHttpUpgradeMap generates the mapping for HTTP Upgrade
func generateHttpUpgradeMap() string {
	return `map {path} {backend_http_upgrade} {
		/trojan 127.0.0.1:8006
		/vmess 127.0.0.1:8005
		/vless 127.0.0.1:8004
	}`
}

// WriteToFile writes the generated Caddyfile to the specified file
func WriteToFile(content, filePath string) error {
	if err := ensureDirectory(filePath); err != nil {
		return fmt.Errorf("failed to ensure directory: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write Caddyfile to %s: %w", filePath, err)
	}

	logger.GetLogger().Info().Str("filepath", filePath).Msg("Caddyfile written successfully")
	return nil
}

// ensureDirectory checks if the directory exists, and creates it if it does not
func ensureDirectory(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}