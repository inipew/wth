package singbox

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type Config struct {
    Log         Log				`json:"log,omitempty"`
    DNS         DNSConfig		`json:"dns,omitempty"`
    NTP         NTP				`json:"ntp,omitempty"`
    Inbounds    []Inbound		`json:"inbounds,omitempty"`
    Outbounds   []Outbound		`json:"outbounds,omitempty"`
    Route       Route			`json:"route,omitempty"`
    Experimental Experimental	`json:"experimental,omitempty"`
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
    Rules        []Rule   `json:"rules,omitempty"`
    ClientSubnet string   `json:"client_subnet,omitempty"`
    // RuleSet			[]string   `json:"rule_set,omitempty"`
    // Outbounds        []string `json:"outbound,omitempty"`
}

type NTP struct {
    Interval    string `json:"interval,omitempty"`
    Server      string `json:"server,omitempty"`
    ServerPort  int    `json:"server_port,omitempty"`
    Detour      string `json:"detour,omitempty"`
}

type User struct {
    Name     string `json:"name,omitempty"`
    UUID     string `json:"uuid"`
    Password string `json:"password"`
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
    Multiplex     Multiplex  `json:"multiplex,omitempty"`
    Transport     Transport  `json:"transport,omitempty"`
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
    Rules      []Rule   `json:"rules,omitempty"`
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

func BuildDNSConfig() DNSConfig {
	return DNSConfig{
		Servers: []DNSServer{
			{
				Tag:             "remote_dns",
				Address:         "https://cloudflare-dns.com/dns-query",
				AddressResolver: "dns_local",
				Strategy:        "prefer_ipv4",
				Detour:          "direct",
			},
			{
				Tag:      "dns_local",
				Address:  "local",
				Strategy: "prefer_ipv4",
				Detour:   "direct",
			},
			{
				Tag:    "dns_block",
				Address: "rcode://success",
			},
		},
		Rules: []DNSRule{
			{
				// RuleSet:      []string{"geosite-malicious", "geoip-malicious"},
				Rules: []Rule{
					{RuleSet: []string{"geosite-malicious", "geoip-malicious"}},
				},
				Server:       "dns_block",
				DisableCache: true,
			},
			{
				Type:       "logical",
				Mode:       "and",
				Rules:      []Rule{
					{Protocol: "quic"}, 
					{RuleSet: []string{"youtube"}},
				},
				Server:     "dns_block",
				DisableCache: true,
				RewriteTTL:  10,
			},
			{
				Rules: []Rule{
					{Outbounds: []string{"any"}},
				},
				Server:       "remote_dns",
				ClientSubnet: "103.3.60.0/22",
			},
		},
		Final:            "remote_dns",
		IndependentCache: true,
	}
}

func AddUserToInbounds(inbounds *[]Inbound, inboundType, transportType string, newUser User) {
	for i := range *inbounds {
		inbound := &(*inbounds)[i]
		if inbound.Type == inboundType {
			if transportType == "" || (inbound.Transport.Type == transportType) {
				inbound.Users = append(inbound.Users, newUser)
			}
		}
	}
}

func BuildInboundConfig(types string, tag string, listenPort int, tcpFastOpen bool, domainStrategy, listen string, sniff bool, sniffTimeout string, users []User, transporttype string, path string) Inbound {
	inbound := Inbound{
		Type:          types,
		Tag:           tag,
		Listen:        listen,
		ListenPort:    listenPort,
		TCPFastOpen:   tcpFastOpen,
		Sniff:         sniff,
		SniffOverrideDestination: false,
		SniffTimeout:  sniffTimeout,
		DomainStrategy: domainStrategy,
		Users:         users,
	}

	// Check if the type is "socks"
	if types != "socks" {
		inbound.Multiplex = Multiplex{Enabled: true}

		inbound.Transport = Transport{
			Type: transporttype,
		}

		// Set Transport properties based on transport type
		switch transporttype {
		case "ws":
			inbound.Transport.Path = "/" + path
			inbound.Transport.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
		case "httpupgrade":
			inbound.Transport.Path = "/" + path
		case "grpc":
			inbound.Transport.ServiceName = path
		}
	} else {
		inbound.Multiplex = Multiplex{}
		inbound.Transport = Transport{}
	}

	return inbound
}


func BuildOutboundConfig(tag, outboundType string, outbounds []string, defaultOutbound string) Outbound {
	config := Outbound{
		Type: outboundType,
		Tag:  tag,
	}

	// Only set specific fields for selector type
	if outboundType == "selector" {
		config.Outbounds = outbounds
		config.Default = defaultOutbound
		config.InterruptExistConnections = true
	}

	return config
}

func BuildExperimentalConf() Experimental{
	return Experimental{
		CacheFile: CacheFile{
				Enabled:   true,
				Path:      "caches.db",
				CacheID:   "sing",
				StoreRDrc: true,
			},
			ClashAPI: ClashAPI{
				ExternalController:      "[::]:9090",
				ExternalUI:             "dashboard",
				ExternalUIDownloadURL:  "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
				ExternalUIDownloadDetour: "direct",
				Secret:                 "qwe12345",
			},
	}
}

func BuildNTPConfig() Config{
    return Config{
        NTP: NTP{
            Interval:   "5m0s",
			Server:     "time.apple.com",
			ServerPort: 123,
			Detour:     "direct",
        },
    }
}

func BuildWireGuardOutboundConfig(ipv6 string, privatekey string) WireGuardOutbound {
	return WireGuardOutbound{
		Type:           "wireguard",
		Tag:            "warp-out",
		DomainStrategy: "prefer_ipv4",
		LocalAddress:   []string{"172.16.0.2/32", ipv6},
		PrivateKey:     privatekey,
		Server:         "engage.cloudflareclient.com",
		ServerPort:     2408,
		PeerPublicKey:  "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
		MTU:            1280,
	}
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

func generateUUID() string {
	return uuid.NewString()
}

