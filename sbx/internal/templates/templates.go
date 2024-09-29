package templates

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