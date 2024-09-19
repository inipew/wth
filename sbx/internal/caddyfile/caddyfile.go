package caddyfile

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sbx/internal/config"
	"text/template"

	"github.com/rs/zerolog/log"
)

// Config holds the information needed to generate the Caddyfile
type Config struct {
	Domain string
}

// CaddyfileTemplate defines the template for the Caddyfile
const CaddyfileTemplate = `
{
	log {
		output file {{ .CaddyLogFile }} {
			roll_keep 15
			roll_keep_for 48h
		}
		format console
	}
	servers {
		trusted_proxies static {{ .TrustedProxies }}
	}
}

{{ .Domain }}, http://{{ .Domain }} {
	@websocket {
		header Connection *Upgrade*
		header Upgrade websocket
		header Sec-WebSocket-Key *
	}
	@http_upgrade {
		header Connection *Upgrade*
		header Upgrade websocket
		not header Sec-WebSocket-Key *
	}
	@grpc {
		header Content-Type "application/grpc"
		protocol grpc
	}

	# Mapping for WebSocket and HTTP Upgrade
	{{ .WebSocketMap }}
	{{ .HttpUpgradeMap }}

	handle @websocket {
		@rewrite_path_websocket {
			path_regexp ^/.*?/(trojan|vmess|vless)
		}
		handle @rewrite_path_websocket {
			rewrite * /{http.regexp.1}
		}
		reverse_proxy {backend}
	}

	handle @http_upgrade {
		@rewrite_path_http_upgrade {
			path_regexp ^/.*?/(trojan|vmess|vless)
		}
		handle @rewrite_path_http_upgrade {
			rewrite * /{http.regexp.1}
		}
		reverse_proxy {backend_http_upgrade}
	}

	handle @grpc {
		reverse_proxy {backend} {
			transport http {
				versions h2c
			}
		}
	}

	header {
		Cache-Control "public, max-age=3600"
		X-Content-Type-Options "nosniff"
		X-Frame-Options "DENY"
		X-XSS-Protection "1; mode=block"
	}
}`

// GenerateCaddyfile generates a Caddyfile based on the provided config
func GenerateCaddyfile(config Config) (string, error) {
	if err := validateConfig(config); err != nil {
		return "", fmt.Errorf("invalid config: %w", err)
	}

	tmpl, err := template.New("caddyfile").Parse(CaddyfileTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, prepareTemplateData(config)); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

// validateConfig checks if the domain is valid
func validateConfig(config Config) error {
	if config.Domain == "" {
		return errors.New("domain cannot be empty")
	}
	return nil
}

// prepareTemplateData prepares the data for the template
func prepareTemplateData(cfg Config) map[string]interface{} {
	return map[string]interface{}{
		"Domain":          cfg.Domain,
		"CaddyLogFile":    config.CaddyLogFile,
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

	log.Info().Str("filepath", filePath).Msg("Caddyfile written successfully")
	return nil
}

// ensureDirectory checks if the directory exists, and creates it if it does not
func ensureDirectory(filepath string) error {
	dir := filepath[:len(filepath)-len("Caddyfile")]
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}