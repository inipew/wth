package main

import (
	"fmt"
	"log"
	"singconf/config/sing"

	"github.com/google/uuid"
)

func main() {
	listen := "0.0.0.0"
	domainStrategy := "prefer_ipv4"
	tcpFastOpen := true
	sniffTimeout := "300ms"
	sniff := true

	// config := sing.Config{
	// 	Log: sing.Log{
	// 		Level:     "info",
	// 		Output:    "/usr/local/sing-box/sing-box.log",
	// 		Timestamp: true,
	// 	},
	// 	DNS: buildDNSConfig(),
	// 	NTP: sing.NTP{
	// 		Interval:   "5m0s",
	// 		Server:     "time.apple.com",
	// 		ServerPort: 123,
	// 		Detour:     "direct",
	// 	},
	// 	Inbounds: []sing.Inbound{
	// 		buildInboundConfig("vless", "vless-ws-in", 8001, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "ws", "vless"),
	// 		buildInboundConfig("vmess", "vmess-ws-in", 8002, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "ws", "vmess"),
	// 		buildInboundConfig("trojan", "trojan-ws-in", 8003, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "ws", "trojan"),
	// 		buildInboundConfig("vless", "vless-httpupgrade-in", 8004, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "httpupgrade", "vless"),
	// 		buildInboundConfig("vmess", "vmess-httpupgrade-in", 8005, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "httpupgrade", "vmess"),
	// 		buildInboundConfig("trojan", "trojan-httpupgrade-in", 8006, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "httpupgrade", "trojan"),
	// 		buildInboundConfig("vless", "vless-grpc-in", 8007, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "grpc", "vless"),
	// 		buildInboundConfig("vmess", "vmess-grpc-in", 8008, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "grpc", "vmess"),
	// 		buildInboundConfig("trojan", "trojan-grpc-in", 8009, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "grpc", "trojan"),
	// 		buildInboundConfig("socks","socks-in",8093,true,"prefer_ipv4", "0.0.0.0", true, "300ms",[]sing.User{
	// 			{Username: "default",Password: generateUUID()},
	// 		},"",""),
	// 	},
	// 	Outbounds: []sing.Outbound{
	// 		buildOutboundConfig("direct", "direct", nil, ""),
	// 		buildOutboundConfig("block", "block", nil, ""),
	// 		buildOutboundConfig("dns-out", "dns", nil, ""),
	// 		buildOutboundConfig("TrafficUmum", "selector", []string{"direct"}, "direct"),
	// 		buildOutboundConfig("TrafficGoogle", "selector", []string{"direct"}, "direct"),
	// 		buildOutboundConfig("TrafficAds", "selector", []string{"direct", "block"}, "block"),
	// 		buildOutboundConfig("TrafficPorn", "selector", []string{"direct", "block"}, "block"),
	// 	},
	// 	Experimental: sing.Experimental{
	// 		CacheFile: sing.CacheFile{
	// 			Enabled:    true,
	// 			Path:       "cache.db",
	// 			StoreRDrc: true,
	// 			CacheID:    "sing-box",
	// 		},
	// 		ClashAPI: sing.ClashAPI{
	// 			ExternalController:        "0.0.0.0:9090",
	// 			ExternalUI:                "yacd",
	// 			ExternalUIDownloadURL:     "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
	// 			ExternalUIDownloadDetour:  "direct",
	// 		},
	// 	},
	// }

	// // Add users to inbounds
	// addUserToInbounds(&config.Inbounds, "vmess", "ws", sing.User{Name: "pew", UUID: generateUUID()})
	// addUserToInbounds(&config.Inbounds, "trojan", "", sing.User{Name: "pew", Password: generateUUID()})

	// // Encode to JSON
	// configJSON, err := json.MarshalIndent(config, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error encoding JSON: %v", err)
	// }

	// // Output JSON to file
	// err = os.WriteFile("config.json", configJSON, 0644)
	// if err != nil {
	// 	log.Fatalf("Error writing JSON to file: %v", err)
	// }
	logDNSExperimental := sing.Config{
		Log: &sing.Log{
			Level:     "info",
			Output:    "/usr/local/sing-box/sing-box.log",
			Timestamp: true,
		},
		DNS: buildDNSConfig(),
		Experimental: BuildExperimentalConf(),
		NTP: BuildNTPConfig(),
	}
	inboundConfig := sing.Config{
		Inbounds: []sing.Inbound{
			buildInboundConfig("vless", "vless-ws-in", 8001, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "ws", "vless"),
			buildInboundConfig("vmess", "vmess-ws-in", 8002, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "ws", "vmess"),
			buildInboundConfig("trojan", "trojan-ws-in", 8003, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", Password: generateUUID()},
			}, "ws", "trojan"),
			buildInboundConfig("vless", "vless-httpupgrade-in", 8004, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "httpupgrade", "vless"),
			buildInboundConfig("vmess", "vmess-httpupgrade-in", 8005, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "httpupgrade", "vmess"),
			buildInboundConfig("trojan", "trojan-httpupgrade-in", 8006, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", Password: generateUUID()},
			}, "httpupgrade", "trojan"),
			buildInboundConfig("vless", "vless-grpc-in", 8007, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "grpc", "vless"),
			buildInboundConfig("vmess", "vmess-grpc-in", 8008, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", UUID: generateUUID()},
			}, "grpc", "vmess"),
			buildInboundConfig("trojan", "trojan-grpc-in", 8009, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []sing.User{
				{Name: "default", Password: generateUUID()},
			}, "grpc", "trojan"),
			buildInboundConfig("socks","socks-in",8093,true,"prefer_ipv4", "0.0.0.0", true, "300ms",[]sing.User{
				{Username: "default",Password: generateUUID()},
			},"",""),
		},
		Route: sing.NewRouteConfig(),
	}
	if err := logDNSExperimental.SaveToFile("config.json"); err != nil{
		log.Fatalf("error save json file: %s", err.Error())
	}
	if err := inboundConfig.SaveToFile("in.json"); err != nil{
		log.Fatalf("error save json file: %s", err.Error())
	}

	fmt.Println("Configuration JSON has been generated successfully.")
}

func generateUUID() string {
	return uuid.NewString()
}

func buildDNSConfig() *sing.DNSConfig {
	return &sing.DNSConfig{
		Servers: []sing.DNSServer{
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
		Rules: []sing.DNSRule{
			{
				RuleSet: []string{"geosite-malicious", "geoip-malicious"},
				Server:       "dns_block",
				DisableCache: true,
			},
			{
				Type:       "logical",
				Mode:       "and",
				Rule:      []sing.Rule{{Protocol: "quic"}, {RuleSet: []string{"youtube"}}},
				Server:     "dns_block",
				DisableCache: true,
				RewriteTTL:  10,
			},
			{
				Outbounds: []string{"any"},
				Server:       "remote_dns",
				ClientSubnet: "103.3.60.0/22",
			},
		},
		Final:            "remote_dns",
		IndependentCache: true,
	}
}

func buildInboundConfig(types string, tag string, listenPort int, tcpFastOpen bool, domainStrategy, listen string, sniff bool, sniffTimeout string, users []sing.User, transporttype string, path string) sing.Inbound {
	inbound := sing.Inbound{
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
		inbound.Multiplex = &sing.Multiplex{Enabled: true}

		inbound.Transport = &sing.Transport{
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

func buildOutboundConfig(tag, outboundType string, outbounds []string, defaultOutbound string) sing.Outbound {
	config := sing.Outbound{
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

func addUserToInbounds(inbounds *[]sing.Inbound, inboundType, transportType string, newUser sing.User) {
	for i := range *inbounds {
		inbound := &(*inbounds)[i]
		if inbound.Type == inboundType {
			if transportType == "" || (inbound.Transport.Type == transportType) {
				inbound.Users = append(inbound.Users, newUser)
			}
		}
	}
}

func BuildExperimentalConf() *sing.Experimental{
	return &sing.Experimental{
		CacheFile: sing.CacheFile{
				Enabled:   true,
				Path:      "caches.db",
				CacheID:   "sing",
				StoreRDrc: true,
			},
			ClashAPI: sing.ClashAPI{
				ExternalController:      "[::]:9090",
				ExternalUI:             "dashboard",
				ExternalUIDownloadURL:  "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
				ExternalUIDownloadDetour: "direct",
				Secret:                 "qwe12345",
			},
	}
}

func BuildNTPConfig() *sing.NTP {
    return &sing.NTP{
		Interval:   "5m0s",
		Server:     "time.apple.com",
		ServerPort: 123,
		Detour:     "direct",
	}
}

func BuildWireGuardOutboundConfig(ipv6 string, privatekey string) sing.WireGuardOutbound {
	return sing.WireGuardOutbound{
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

