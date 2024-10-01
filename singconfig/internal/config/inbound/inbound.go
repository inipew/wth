package inbound

import (
	"singconfig/internal/utils"
)

type InboundConfig struct {
	Type                    string              `json:"type"`
	Tag                     string              `json:"tag"`
	Listen                  string              `json:"listen"`
	ListenPort              int                 `json:"listen_port"`
	TCPFastOpen             bool                `json:"tcp_fast_open,omitempty"`
	Sniff                   bool                `json:"sniff,omitempty"`
	SniffTimeout            string              `json:"sniff_timeout,omitempty"`
	SniffOverrideDestination bool                `json:"sniff_override_destination,omitempty"`
	DomainStrategy          string              `json:"domain_strategy,omitempty"`
	Users                   []UserConfig        `json:"users,omitempty"`
	Multiplex               *MultiplexConfig     `json:"multiplex,omitempty"`
	Transport               *TransportConfig     `json:"transport,omitempty"`
}

type UserConfig struct {
	Name     string `json:"name,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Password string `json:"password,omitempty"`
	Username string `json:"Username,omitempty"`
}

type MultiplexConfig struct {
	Enabled bool `json:"enabled"`
}

type TransportConfig struct {
	Type                 string `json:"type"`
	Path                 string `json:"path,omitempty"`
	EarlyDataHeaderName  string `json:"early_data_header_name,omitempty"`
	ServiceName          string `json:"service_name,omitempty"`
}

func BuildInbounds() []InboundConfig {
	inboundTypes := []struct {
		protocol    string
		tag         string
		listenPort  int
		tcpFastOpen bool
		transport   string
		path        string
	}{
		{"vless", "vless-ws-in", 8001, true, "ws", "vless"},
		{"vmess", "vmess-ws-in", 8002, true, "ws", "vmess"},
		{"trojan", "trojan-ws-in", 8003, true, "ws", "trojan"},
		{"vless", "vless-httpupgrade-in", 8004, true, "httpupgrade", "vless"},
		{"vmess", "vmess-httpupgrade-in", 8005, true, "httpupgrade", "vmess"},
		{"trojan", "trojan-httpupgrade-in", 8006, true, "httpupgrade", "trojan"},
		{"vless", "vless-grpc-in", 8007, true, "grpc", "vless"},
		{"vmess", "vmess-grpc-in", 8008, true, "grpc", "vmess"},
		{"trojan", "trojan-grpc-in", 8009, true, "grpc", "trojan"},
		{"socks", "socks-in", 8093, true, "", ""},
	}

	var inbounds []InboundConfig
	uuid := utils.GenerateUUID()
	for _, entry := range inboundTypes {
		inbounds = append(inbounds, buildInboundConfig(entry, uuid))
	}
	return inbounds
}

func buildInboundConfig(entry struct {
	protocol    string
	tag         string
	listenPort  int
	tcpFastOpen bool
	transport   string
	path        string
}, uuid string) InboundConfig {
	users := buildUsers(entry.protocol, uuid)
	inbound := InboundConfig{
		Type:                     entry.protocol,
		Tag:                      entry.tag,
		Listen:                   "0.0.0.0",
		ListenPort:               entry.listenPort,
		TCPFastOpen:              entry.tcpFastOpen,
		Sniff:                    true,
		SniffOverrideDestination: false,
		SniffTimeout:             "300ms",
		DomainStrategy:           "prefer_ipv4",
		Users:                    users,
	}

	if entry.protocol != "socks" {
		inbound.Multiplex = &MultiplexConfig{Enabled: true}
		inbound.Transport = buildTransportConfig(entry.transport, entry.path)
	}

	return inbound
}

func buildUsers(protocol string, uuid string) []UserConfig {
	switch protocol {
	case "trojan":
		return []UserConfig{{Name: "default", Password: uuid}}
	case "socks":
		return []UserConfig{{Username: "default", Password: uuid}}
	default:
		return []UserConfig{{Name: "default", UUID: uuid}}
	}
}

func buildTransportConfig(transportType, path string) *TransportConfig {
	transport := &TransportConfig{
		Type: transportType,
	}

	switch transportType {
	case "ws":
		transport.Path = "/" + path
		transport.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
	case "httpupgrade":
		transport.Path = "/" + path
	case "grpc":
		transport.ServiceName = path
	}

	return transport
}