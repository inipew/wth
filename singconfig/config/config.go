package config

import (
	"encoding/json"
	"os"
)

type SingBoxConfig struct {
	Log         *LogConfig         `json:"log,omitempty"`
	DNS         *DNSConfig         `json:"dns,omitempty"`
	NTP         *NTPConfig         `json:"ntp,omitempty"`
	Inbounds    []InboundConfig   `json:"inbounds,omitempty"`
	Outbounds   []OutboundConfig  `json:"outbounds,omitempty"`
	Route       *RouteConfig       `json:"route,omitempty"`
	Experimental *ExperimentalConfig `json:"experimental,omitempty"`
}

type LogConfig struct {
	Disabled  bool	`json:"disabled"`
	Level     string `json:"level"`
	Output    string `json:"output"`
	Timestamp bool   `json:"timestamp"`
}

type DNSConfig struct {
	Servers          []DNSServerConfig `json:"servers"`
	Rules            []DNSRuleConfig   `json:"rules"`
	Final            string            `json:"final"`
	IndependentCache bool              `json:"independent_cache"`
}

type DNSServerConfig struct {
	Tag              string `json:"tag"`
	Address          string `json:"address"`
	AddressResolver  string `json:"address_resolver,omitempty"`
	Strategy         string `json:"strategy,omitempty"`
	Detour           string `json:"detour,omitempty"`
}

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

type NTPConfig struct {
	Interval    string `json:"interval"`
	Server      string `json:"server"`
	ServerPort  int    `json:"server_port"`
	Detour      string `json:"detour"`
}

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

type WireGuardOutbound struct{
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

type RouteConfig struct {
	Rules    []RouteRuleConfig `json:"rules,omitempty"`
	RuleSet  []RuleSetConfig   `json:"rule_set,omitempty"`
	Final    string            `json:"final,omitempty"`
	AutoDetectInterface bool   `json:"auto_detect_interface,omitempty"`
}

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

type RuleSetConfig struct {
	Type            string        `json:"type"`
	Tag             string        `json:"tag"`
	Format          string        `json:"format"`
	URL             string        `json:"url"`
	DownloadDetour  string        `json:"download_detour"`
	UpdateInterval  string        `json:"update_interval"`
}

type ExperimentalConfig struct {
	CacheFile CacheFileConfig `json:"cache_file"`
	ClashAPI  ClashAPIConfig  `json:"clash_api"`
}

type CacheFileConfig struct {
	Enabled   bool   `json:"enabled"`
	Path      string `json:"path"`
	CacheID   string `json:"cache_id"`
	StoreRDRC bool   `json:"store_rdrc"`
}

type ClashAPIConfig struct {
	ExternalController     string `json:"external_controller"`
	ExternalUI             string `json:"external_ui"`
	ExternalUIDownloadURL  string `json:"external_ui_download_url"`
	ExternalUIDownloadDetour string `json:"external_ui_download_detour"`
	Secret                 string `json:"secret"`
}

func (c *SingBoxConfig) SaveToFile(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(c)
}