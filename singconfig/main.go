package main

import (
	"log"
	"singconfig/config"

	"github.com/google/uuid"
)

func main(){
	conf := config.SingBoxConfig{
		Log: buildLogConfig(),
		DNS: buildDNSConfig(),
		NTP: buildNTPConfig(),
		Experimental: buildExperimentalConfig(),
		Inbounds: buildInbounds(),
		Outbounds: 	buildOutbounds(),
		Route: NewRouteConfig(),
	}

	if err := conf.SaveToFile("config.json"); err != nil{
		log.Fatalf("error save json file: %s", err.Error())
	}
}

func generateUUID() string {
	return uuid.NewString()
}

func buildLogConfig() *config.LogConfig {
	return &config.LogConfig{
		Disabled:  false,
		Level:     "info",
		Output:    "/usr/local/sing-box/sing-box.log",
		Timestamp: true,
	}
}

func buildNTPConfig() *config.NTPConfig {
	return &config.NTPConfig{
		Interval:   "5m0s",
		Server:     "time.apple.com",
		ServerPort: 123,
		Detour:     "direct",
	}
}

func buildExperimentalConfig() *config.ExperimentalConfig {
	return &config.ExperimentalConfig{
		CacheFile: config.CacheFileConfig{
			Enabled:   true,
			Path:      "caches.db",
			CacheID:   "sing",
			StoreRDRC: true,
		},
		ClashAPI: config.ClashAPIConfig{
			ExternalController:      "[::]:9090",
			ExternalUI:             "dashboard",
			ExternalUIDownloadURL:  "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
			ExternalUIDownloadDetour: "direct",
			Secret:                 "qwe12345",
		},
	}
}

func buildDNSConfig() *config.DNSConfig {
	return &config.DNSConfig{
		Servers: []config.DNSServerConfig{
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
		Rules: []config.DNSRuleConfig{
			{
				RuleSet: []string{"geosite-malicious", "geoip-malicious"},
				Server:       "dns_block",
				DisableCache: true,
			},
			{
				Type:       "logical",
				Mode:       "and",
				Rules:      []config.DNSRuleConfig{{Protocol: "quic"}, {RuleSet: []string{"youtube"}}},
				Server:     "dns_block",
				DisableCache: true,
				RewriteTTL:  10,
			},
			{
				Outbound: []string{"any"},
				Server:       "remote_dns",
				ClientSubnet: "103.3.60.0/22",
			},
		},
		Final:            "remote_dns",
		IndependentCache: true,
	}
}

func buildInbounds() []config.InboundConfig {
	inboundTypes := []struct {
		protocol     string
		tag          string
		listenPort   int
		tcpFastOpen  bool
		transport    string
		path		 string
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

	var inbounds []config.InboundConfig
	for _, entry := range inboundTypes {
		users := []config.UserConfig{{Name: "default", UUID: generateUUID()}}
		if entry.protocol == "trojan" {
			users = []config.UserConfig{{Name: "default", Password: generateUUID()}}
		} else if entry.protocol == "socks" {
			users = []config.UserConfig{{Username: "default", Password: generateUUID()}}
		}
		inbounds = append(inbounds, buildInboundConfig(entry.protocol, entry.tag, entry.listenPort, entry.tcpFastOpen, "prefer_ipv4", "0.0.0.0", true, "300ms", users, entry.transport, entry.path))
	}
	return inbounds
}

func buildInboundConfig(types string, tag string, listenPort int, tcpFastOpen bool, domainStrategy, listen string, sniff bool, sniffTimeout string, users []config.UserConfig, transporttype string, path string) config.InboundConfig {
	inbound := config.InboundConfig{
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
		inbound.Multiplex = &config.MultiplexConfig{Enabled: true}

		inbound.Transport = &config.TransportConfig{
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
	}

	return inbound
}

func buildOutbounds() []config.OutboundConfig {
	return []config.OutboundConfig{
		buildOutboundConfig("direct", "direct", nil, ""),
		buildOutboundConfig("block", "block", nil, ""),
		buildOutboundConfig("dns-out", "dns", nil, ""),
		buildOutboundConfig("TrafficUmum", "selector", []string{"direct"}, "direct"),
		buildOutboundConfig("TrafficGoogle", "selector", []string{"direct"}, "direct"),
		buildOutboundConfig("TrafficAds", "selector", []string{"direct", "block"}, "block"),
		buildOutboundConfig("TrafficPorn", "selector", []string{"direct", "block"}, "block"),
	}
}

func buildOutboundConfig(tag, outboundType string, outbounds []string, defaultOutbound string) config.OutboundConfig {
	config := config.OutboundConfig{
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

func NewRouteConfig() *config.RouteConfig {
	return &config.RouteConfig{
		Rules:               createRules(),
		RuleSet:             createRuleSets(),
		Final:               "TrafficUmum",
		AutoDetectInterface: true,
	}
}

func createRules() []config.RouteRuleConfig {
	return []config.RouteRuleConfig{
		createDNSRule(),
		createDirectRule(),
		createMaliciousRule(),
		createQuicYoutubeRule(),
		createAdsRule(),
		createPornRule(),
		createCloudflareRule(),
		createGoogleRule(),
		createWarpRule(),
		createCommonPortsRule(),
	}
}

func createDNSRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		Type: "logical",
		Mode: "or",
		Rules: []config.RouteRuleConfig{
			{Port: 53},
			{Protocol: "dns"},
		},
		Outbound: "dns-out",
	}
}

func createDirectRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"direct_some_web", "geosite-google_ads"},
		Outbound: "direct",
	}
}

func createMaliciousRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"geosite-malicious", "geoip-malicious"},
		Outbound: "block",
	}
}

func createQuicYoutubeRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		Type: "logical",
		Mode: "and",
		Rules: []config.RouteRuleConfig{
			{Protocol: "quic"},
			{RuleSet: []string{"rule_youtube"}},
		},
		Outbound: "block",
	}
}

func createAdsRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"oisd-full", "rule-ads", "d3ward"},
		Outbound: "TrafficAds",
	}
}

func createPornRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"oisd-nsfw", "category-porn"},
		Outbound: "TrafficPorn",
	}
}

func createCloudflareRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		DomainSuffix: []string{"gstatic.com", "cp.cloudflare.com"},
		Outbound:    "TrafficUmum",
	}
}

func createGoogleRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"geoip_google", "google", "rule_youtube"},
		Outbound: "TrafficGoogle",
	}
}

func createWarpRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"openai", "needwarp", "reddit", "geoip_cloudflare", "microsoft", "rule_github"},
		Outbound: "TrafficWarp",
	}
}

func createCommonPortsRule() config.RouteRuleConfig {
	return config.RouteRuleConfig{
		RuleSet:   []string{"googlefcm", "commonports"},
		Outbound: "TrafficUmum",
	}
}

func createRuleSets() []config.RuleSetConfig {
	return []config.RuleSetConfig{
		createRemoteRuleSet("direct_some_web", "https://cdn.jsdelivr.net/gh/inipew/any@main/direct-some-web.srs", "24h0m0s"),
		createRemoteRuleSet("needwarp", "https://cdn.jsdelivr.net/gh/inipew/any@main/warped.srs", "24h0m0s"),
		createRemoteRuleSet("adguard", "https://cdn.jsdelivr.net/gh/inipew/any@main/adguard.srs", "24h0m0s"),
		createRemoteRuleSet("geosite-google_ads", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-google-ads.srs", "24h0m0s"),
		createRemoteRuleSet("geoip-malicious", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geoip/geoip-malicious.srs", "24h0m0s"),
		createRemoteRuleSet("oisd-full", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-oisd-full.srs", "168h0m0s"),
		createRemoteRuleSet("oisd-nsfw", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-oisd-nsfw.srs", "168h0m0s"),
		createRemoteRuleSet("rule-ads", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-rule-ads.srs", "168h0m0s"),
		createRemoteRuleSet("d3ward", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-d3ward.srs", "168h0m0s"),
		createRemoteRuleSet("geosite-malicious", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-rule-malicious.srs", "168h0m0s"),
		createRemoteRuleSet("category-porn", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-category-porn.srs", "168h0m0s"),
		createRemoteRuleSet("openai", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-openai.srs", "168h0m0s"),
		createRemoteRuleSet("onedrive", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-onedrive.srs", "168h0m0s"),
		createRemoteRuleSet("microsoft", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-microsoft.srs", "168h0m0s"),
		createRemoteRuleSet("rule_github", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-github.srs", "168h0m0s"),
		createRemoteRuleSet("rule_youtube", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-youtube.srs", "168h0m0s"),
		createRemoteRuleSet("google", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-google.srs", "168h0m0s"),
		createRemoteRuleSet("googlefcm", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-googlefcm.srs", "168h0m0s"),
		createRemoteRuleSet("reddit", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-reddit.srs", "168h0m0s"),
		createRemoteRuleSet("geoip_google", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geoip/geoip-google.srs", "168h0m0s"),
		createRemoteRuleSet("geoip_cloudflare", "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geoip/geoip-cloudflare.srs", "168h0m0s"),
		createRemoteRuleSet("commonports", "https://cdn.jsdelivr.net/gh/inipew/any@main/commonports.srs", "720h0m0s"),
	}
}

func createRemoteRuleSet(tag, url, updateInterval string) config.RuleSetConfig {
	return config.RuleSetConfig{
		Type:           "remote",
		Tag:            tag,
		Format:         "binary",
		URL:            url,
		DownloadDetour: "direct",
		UpdateInterval: updateInterval,
	}
}