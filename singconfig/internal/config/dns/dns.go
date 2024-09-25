package dns

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

func BuildDNSConfig() *DNSConfig {
	return &DNSConfig{
		Servers:          buildDNSServers(),
		Rules:            buildDNSRules(),
		Final:            "remote_dns",
		IndependentCache: true,
	}
}

func buildDNSServers() []DNSServerConfig {
	return []DNSServerConfig{
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
			Tag:     "dns_block",
			Address: "rcode://success",
		},
	}
}

func buildDNSRules() []DNSRuleConfig {
	return []DNSRuleConfig{
		{
			RuleSet:      []string{"geosite-malicious", "geoip-malicious"},
			Server:       "dns_block",
			DisableCache: true,
		},
		{
			Type:         "logical",
			Mode:         "and",
			Rules:        []DNSRuleConfig{{Protocol: "quic"}, {RuleSet: []string{"rule_youtube"}}},
			Server:       "dns_block",
			DisableCache: true,
			RewriteTTL:   10,
		},
		{
			Outbound:     []string{"any"},
			Server:       "remote_dns",
			ClientSubnet: "103.3.60.0/22",
		},
	}
}