package sing

import (
	"encoding/json"
	"os"
)

type Config struct {
    Log         *Log				`json:"log,omitempty"`
    DNS         *DNSConfig		`json:"dns,omitempty"`
    NTP         *NTP            `json:"ntp,omitempty"`
    Inbounds    []Inbound		`json:"inbounds,omitempty"`
    Outbounds   []Outbound		`json:"outbounds,omitempty"`
    Route       *Route          `json:"route,omitempty"`
    Experimental *Experimental	`json:"experimental,omitempty"`
}

type Log struct {
    Level     string `json:"level,omitempty"`
    Output    string `json:"output,omitempty"`
    Timestamp bool   `json:"timestamp,omitempty"`
}

type DNSConfig struct {
	Servers          []DNSServer `json:"servers"`
	Rules            []DNSRule   `json:"rules"`
	Final            string      `json:"final"`
	IndependentCache bool        `json:"independent_cache,omitempty"`
}

type DNSServer struct {
    Tag             string `json:"tag"`
    Address         string `json:"address"`
    AddressResolver string `json:"address_resolver,omitempty"`
    Strategy        string `json:"strategy,omitempty"`
    Detour          string `json:"detour,omitempty"`
}

type DNSRule struct {
    Server       string   `json:"server,omitempty"`
    DisableCache bool     `json:"disable_cache,omitempty"`
    RewriteTTL   int      `json:"rewrite_ttl,omitempty"`
    Type         string   `json:"type,omitempty"`
    Mode         string   `json:"mode,omitempty"`
    Rule        []Rule   `json:"rules,omitempty"`
    ClientSubnet string   `json:"client_subnet,omitempty"`
    RuleSet     []string   `json:"rule_set,omitempty"`
    Outbounds        []string `json:"outbound,omitempty"`
}

type NTP struct {
    Interval    string `json:"interval,omitempty"`
    Server      string `json:"server,omitempty"`
    ServerPort  int    `json:"server_port,omitempty"`
    Detour      string `json:"detour,omitempty"`
}

type User struct {
    Name     string `json:"name,omitempty"`
    UUID     string `json:"uuid,omitempty"`
    Password string `json:"password,omitempty"`
    Username string `json:"Username,omitempty"`
}

type Transport struct {
    Type                   string `json:"type,omitempty"`
    Path                   string `json:"path,omitempty"`
    EarlyDataHeaderName    string `json:"early_data_header_name,omitempty"`
    ServiceName            string `json:"service_name,omitempty"`
}

type Multiplex struct {
    Enabled bool `json:"enabled,omitempty"`
}

type Inbound struct {
    Type          string     `json:"type"`
    Tag           string     `json:"tag"`
    Listen        string     `json:"listen"`
    ListenPort    int        `json:"listen_port"`
    TCPFastOpen   bool       `json:"tcp_fast_open,omitempty"`
    Sniff         bool       `json:"sniff,omitempty"`
    SniffOverrideDestination bool `json:"sniff_override_destination"`
    SniffTimeout  string     `json:"sniff_timeout,omitempty"`
    DomainStrategy string    `json:"domain_strategy,omitempty"`
    Users         []User     `json:"users"`
    Multiplex     *Multiplex  `json:"multiplex,omitempty"`
    Transport     *Transport  `json:"transport,omitempty"`
}

type Outbound struct {
    Type             string   `json:"type"`
    Tag              string   `json:"tag"`
    Outbounds        []string `json:"outbounds,omitempty"`
    Default          string   `json:"default,omitempty"`
    InterruptExistConnections bool `json:"interrupt_exist_connections,omitempty"`
}

type Rule struct {
    Protocol    string   `json:"protocol,omitempty"`
    Port        int      `json:"port,omitempty"`
    DomainSuffix []string `json:"domain_suffix,omitempty"`
    RuleSet     []string   `json:"rule_set,omitempty"`
    Outbounds        []string `json:"outbound,omitempty"`
}

type RouteRule struct {
    Type       string   `json:"type,omitempty"`
    Mode       string   `json:"mode,omitempty"`
    Rule            []Rule   `json:"rules,omitempty"`
    RuleSet         []string   `json:"rule_set,omitempty"`
    Outbounds       string `json:"outbound,omitempty"`
    DomainSuffix    []string `json:"domain_suffix,omitempty"`
}

type RuleSet struct {
    Type            string `json:"type,omitempty"`
    Tag             string `json:"tag,omitempty"`
    Format          string `json:"format,omitempty"`
    URL             string `json:"url,omitempty"`
    DownloadDetour  string `json:"download_detour,omitempty"`
    UpdateInterval  string `json:"update_interval,omitempty"`
}

type Route struct {
    Rules      []RouteRule `json:"rules,omitempty"`
    RuleSet    []RuleSet   `json:"rule_set,omitempty"`
    Final      string      `json:"final,omitempty"`
    AutoDetectInterface bool `json:"auto_detect_interface,omitempty"`
}

type Experimental struct {
    CacheFile	CacheFile	`json:"cache_file,omitempty"`
    ClashAPI	ClashAPI	`json:"clash_api,omitempty"`
}

type CacheFile struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Path     string `json:"path,omitempty"`
	CacheID  string `json:"cache_id,omitempty"`
	StoreRDrc bool  `json:"store_rdrc,omitempty"`
}

type ClashAPI struct {
	ExternalController        string `json:"external_controller,omitempty"`
	ExternalUI                string `json:"external_ui,omitempty"`
	ExternalUIDownloadURL     string `json:"external_ui_download_url,omitempty"`
	ExternalUIDownloadDetour  string `json:"external_ui_download_detour,omitempty"`
	Secret                    string `json:"secret,omitempty"`
}

type WireGuardOutbound struct {
	Type            string   `json:"type"`
	Tag             string   `json:"tag"`
	DomainStrategy  string   `json:"domain_strategy"`
	LocalAddress    []string `json:"local_address"`
	PrivateKey      string   `json:"private_key"`
	Server          string   `json:"server"`
	ServerPort      int      `json:"server_port"`
	PeerPublicKey   string   `json:"peer_public_key"`
	MTU             int      `json:"mtu"`
}

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