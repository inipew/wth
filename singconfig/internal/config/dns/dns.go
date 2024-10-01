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
        newDNSServer("remote_dns", "https://cloudflare-dns.com/dns-query", "dns_local", "prefer_ipv4", "direct"),
        newDNSServer("dns_local", "local", "", "prefer_ipv4", "direct"),
        newDNSServer("dns_block", "rcode://success", "", "", ""),
    }
}

func newDNSServer(tag, address, addressResolver, strategy, detour string) DNSServerConfig {
    return DNSServerConfig{
        Tag:             tag,
        Address:         address,
        AddressResolver: addressResolver,
        Strategy:        strategy,
        Detour:          detour,
    }
}

func buildDNSRules() []DNSRuleConfig {
    return []DNSRuleConfig{
        newBlockRule([]string{"geosite-malicious", "geoip-malicious"}),
        newLogicalRule("and", []DNSRuleConfig{
            {Protocol: "quic"},
            {RuleSet: []string{"rule_youtube"}},
        }, "dns_block", true, 10),
        {
            Outbound:     []string{"any"},
            Server:       "remote_dns",
            ClientSubnet: "103.3.60.0/22",
        },
    }
}

func newBlockRule(ruleSet []string) DNSRuleConfig {
    return DNSRuleConfig{
        RuleSet:      ruleSet,
        Server:       "dns_block",
        DisableCache: true,
    }
}

func newLogicalRule(mode string, rules []DNSRuleConfig, server string, disableCache bool, rewriteTTL int) DNSRuleConfig {
    return DNSRuleConfig{
        Type:         "logical",
        Mode:         mode,
        Rules:        rules,
        Server:       server,
        DisableCache: disableCache,
        RewriteTTL:   rewriteTTL,
    }
}

// func (s DNSServerConfig) withAddress(address string) DNSServerConfig{
// 	s.Address = address
// 	return s
// }

// func (s DNSServerConfig) withAddressResolver(addressResolver string) DNSServerConfig{
// 	s.AddressResolver = addressResolver
// 	return s
// }

// func (s DNSServerConfig) withStrategy(strategy string) DNSServerConfig{
// 	s.Strategy = strategy
// 	return s
// }

// func (s DNSServerConfig) withDetour(detour string) DNSServerConfig{
// 	s.Detour = detour
// 	return s
// }

// func (r DNSRuleConfig) withRuleSet(ruleSet []string) DNSRuleConfig {
// 	r.RuleSet = ruleSet
// 	return r
// }

// func (r DNSRuleConfig) withLogicalRule(mode string, rules []DNSRuleConfig, outbound []string) DNSRuleConfig {
// 	r.Type = "logical"
// 	r.Mode = mode
// 	r.Rules = rules
// 	r.Outbound = outbound
// 	return r
// }

// func (r DNSRuleConfig) withServer(server string) DNSRuleConfig {
// 	r.Server = server
// 	return r
// }

// func (r DNSRuleConfig) withClientSubner(client_subnet string) DNSRuleConfig {
// 	r.ClientSubnet = client_subnet
// 	return r
// }

// func (r DNSRuleConfig) withDisableCache(disableCache bool) DNSRuleConfig {
// 	r.DisableCache = disableCache
// 	return r
// }

// func (r DNSRuleConfig) withRewriteTTL(rewriteTTL int) DNSRuleConfig {
// 	r.RewriteTTL = rewriteTTL
// 	return r
// }

// func (r DNSRuleConfig) withProtocol(protocol string) DNSRuleConfig {
// 	r.Protocol = protocol
// 	return r
// }