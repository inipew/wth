package sing

import (
	"encoding/json"
	"os"
)

// Config adalah struktur utama untuk konfigurasi
type Config struct {
	Log          *LogConfig         `json:"log,omitempty"`
	DNS          *DNSConfig         `json:"dns,omitempty"`
	NTP          *NTPConfig         `json:"ntp,omitempty"`
	Inbounds     []InboundConfig    `json:"inbounds,omitempty"`
	Outbounds    []OutboundConfig   `json:"outbounds,omitempty"`
	Route        *RouteConfig       `json:"route,omitempty"`
	Experimental *ExperimentalConfig `json:"experimental,omitempty"`
}

// LogConfig menangani konfigurasi logging
type LogConfig struct {
	Disabled  bool   `json:"disabled"`
	Level     string `json:"level"`
	Output    string `json:"output"`
	Timestamp bool   `json:"timestamp"`
}

// DNSConfig menangani konfigurasi DNS
type DNSConfig struct {
	Servers          []DNSServerConfig `json:"servers"`
	Rules            []DNSRuleConfig   `json:"rules"`
	Final            string            `json:"final"`
	IndependentCache bool              `json:"independent_cache"`
}

// DNSServerConfig adalah konfigurasi untuk server DNS
type DNSServerConfig struct {
	Tag              string `json:"tag"`
	Address          string `json:"address"`
	AddressResolver  string `json:"address_resolver,omitempty"`
	Strategy         string `json:"strategy,omitempty"`
	Detour           string `json:"detour,omitempty"`
}

// DNSRuleConfig adalah konfigurasi untuk aturan DNS
type DNSRuleConfig struct {
	RuleSet        []string          `json:"rule_set,omitempty"`
	Server         string            `json:"server,omitempty"`
	DisableCache   bool              `json:"disable_cache,omitempty"`
	RewriteTTL     int               `json:"rewrite_ttl,omitempty"`
	Type           string            `json:"type,omitempty"`
	Mode           string            `json:"mode,omitempty"`
	Rules          []DNSRuleConfig   `json:"rules,omitempty"`
	Protocol       string            `json:"protocol,omitempty"`
	Outbound       []string          `json:"outbound,omitempty"`
	ClientSubnet   string            `json:"client_subnet,omitempty"`
}

// NTPConfig menangani konfigurasi NTP
type NTPConfig struct {
	Interval    string `json:"interval"`
	Server      string `json:"server"`
	ServerPort  int    `json:"server_port"`
	Detour      string `json:"detour"`
}

// InboundConfig menangani konfigurasi inbound
type InboundConfig struct {
	Type                     string           `json:"type"`
	Tag                      string           `json:"tag"`
	Listen                   string           `json:"listen"`
	ListenPort               int              `json:"listen_port"`
	TCPFastOpen              bool             `json:"tcp_fast_open,omitempty"`
	Sniff                    bool             `json:"sniff,omitempty"`
	SniffTimeout             string           `json:"sniff_timeout,omitempty"`
	SniffOverrideDestination bool             `json:"sniff_override_destination,omitempty"`
	DomainStrategy           string           `json:"domain_strategy,omitempty"`
	Users                    []UserConfig     `json:"users,omitempty"`
	Multiplex                *MultiplexConfig `json:"multiplex,omitempty"`
	Transport                *TransportConfig `json:"transport,omitempty"`
}

// UserConfig adalah konfigurasi untuk pengguna
type UserConfig struct {
	Name     string `json:"name,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Password string ` json:"password,omitempty"`
	Username string `json:"Username,omitempty"`
}

// MultiplexConfig adalah konfigurasi untuk multiplex
type MultiplexConfig struct {
	Enabled bool `json:"enabled"`
}

// TransportConfig adalah konfigurasi untuk transport
type TransportConfig struct {
	Type                 string `json:"type"`
	Path                 string `json:"path,omitempty"`
	EarlyDataHeaderName  string `json:"early_data_header_name,omitempty"`
	ServiceName          string `json:"service_name,omitempty"`
}

// OutboundConfig menangani konfigurasi outbound
type OutboundConfig struct {
	Type                      string     `json:"type"`
	Tag                       string     `json:"tag"`
	Outbounds                 []string   `json:"outbounds,omitempty"`
	Default                   string     `json:"default,omitempty"`
	InterruptExistConnections bool       `json:"interrupt_exist_connections,omitempty"`
	Detour                    string     `json:"detour,omitempty"`
	DomainStrategy            string     `json:"domain_strategy,omitempty"`
	Interval                  string     `json:"interval,omitempty"`
	IdleTimeout               string     `json:"idle_timeout,omitempty"`
}

// RouteConfig menangani konfigurasi routing
type RouteConfig struct {
	Rules    []RouteRuleConfig `json:"rules,omitempty"`
	RuleSet  []RuleSetConfig   `json:"rule_set,omitempty"`
	Final    string            `json:"final,omitempty"`
	AutoDetectInterface bool   `json:"auto_detect_interface,omitempty"`
}

// RouteRuleConfig adalah konfigurasi untuk aturan routing
type RouteRuleConfig struct {
	Type           string            `json:"type,omitempty"`
	Mode           string            `json:"mode,omitempty"`
	Rules          []RouteRuleConfig `json:"rules,omitempty"`
	Outbound       string            `json:"outbound,omitempty"`
	Port           int               `json:"port,omitempty"`
	Protocol       string            `json:"protocol,omitempty"`
	RuleSet        []string          `json:"rule_set,omitempty"`
	DomainSuffix   []string          `json:"domain_suffix,omitempty"`
}

// RuleSetConfig adalah konfigurasi untuk set aturan
type RuleSetConfig struct {
	Type            string        `json:"type"`
	Tag             string        `json:"tag"`
	Format          string        `json:"format"`
	URL             string        `json:"url"`
	DownloadDetour  string        `json:"download_detour"`
	UpdateInterval  string        `json:"update_interval"`
}

// ExperimentalConfig menangani konfigurasi experimental
type ExperimentalConfig struct {
	CacheFile CacheFileConfig `json:"cache_file"`
	ClashAPI  ClashAPIConfig  `json:"clash_api"`
}

// CacheFileConfig adalah konfigurasi untuk cache file
type CacheFileConfig struct {
	Enabled   bool   `json:"enabled"`
	Path      string `json:"path"`
	CacheID   string `json:"cache_id"`
	StoreRDRC bool   `json:"store_rdrc"`
}

// ClashAPIConfig adalah konfigurasi untuk Clash API
type ClashAPIConfig struct {
	ExternalController     string `json:"external_controller"`
	ExternalUI             string `json:"external_ui"`
	ExternalUIDownloadURL  string `json:"external_ui_download_url"`
	ExternalUIDownloadDetour string `json:"external_ui_download_detour"`
	Secret                 string `json:"secret"`
}

// WireGuardOutbound adalah konfigurasi untuk WireGuard outbound
type WireGuardOutbound struct {
	Type            string   `json:"type,omitempty"`
	Tag             string   `json:"tag,omitempty"`
	DomainStrategy            string     `json:"domain_strategy,omitempty"`
	LocalAddress              []string   `json:"local_address,omitempty"`
	PrivateKey                string     `json:"private_key,omitempty"`
	Server                    string     `json:"server,omitempty"`
	ServerPort                int        `json:"server_port,omitempty"`
	PeerPublicKey             string     `json:"peer_public_key,omitempty"`
	MTU                       int        `json:"mtu,omitempty"`
	URL                       string     `json:"url,omitempty"`
}

// SaveToFile menyimpan konfigurasi ke file
func (c *Config) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}