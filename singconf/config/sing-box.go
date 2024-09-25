package singbox

import (
	"log"
	"path/filepath"
)

func CreateSingBoxConfig(){
	// listen := "0.0.0.0"
	// domainStrategy := "prefer_ipv4"
	// tcpFastOpen := true
	// sniffTimeout := "300ms"
	// sniff := true
	// conf := Config{
	// 	Log: Log{
	// 		Level:     "info",
	// 		Output:    config.SingboxLogFilePath,
	// 		Timestamp: true,
	// 	},
	// 	NTP: NTP{
	// 		Interval:   "5m0s",
	// 		Server:     "time.apple.com",
	// 		ServerPort: 123,
	// 		Detour:     "direct",
	// 	},
	// 	Experimental: Experimental{
	// 		CacheFile: CacheFile{
	// 			Enabled:    true,
	// 			Path:       "cache.db",
	// 			StoreRDrc: true,
	// 			CacheID:    "sing-box",
	// 		},
	// 		ClashAPI: ClashAPI{
	// 			ExternalController:        "0.0.0.0:9090",
	// 			ExternalUI:                "yacd",
	// 			ExternalUIDownloadURL:     "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
	// 			ExternalUIDownloadDetour:  "direct",
	// 		},
	// 	},
	// 	DNS: BuildDNSConfig(),
	// 	Inbounds: []Inbound{
	// 		BuildInboundConfig("vless", "vless-ws-in", 8001, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "ws", "vless"),
	// 		BuildInboundConfig("vmess", "vmess-ws-in", 8002, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "ws", "vmess"),
	// 		BuildInboundConfig("trojan", "trojan-ws-in", 8003, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "ws", "trojan"),
	// 		BuildInboundConfig("vless", "vless-httpupgrade-in", 8004, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "httpupgrade", "vless"),
	// 		BuildInboundConfig("vmess", "vmess-httpupgrade-in", 8005, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "httpupgrade", "vmess"),
	// 		BuildInboundConfig("trojan", "trojan-httpupgrade-in", 8006, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "httpupgrade", "trojan"),
	// 		BuildInboundConfig("vless", "vless-grpc-in", 8007, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "grpc", "vless"),
	// 		BuildInboundConfig("vmess", "vmess-grpc-in", 8008, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", UUID: generateUUID()},
	// 		}, "grpc", "vmess"),
	// 		BuildInboundConfig("trojan", "trojan-grpc-in", 8009, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
	// 			{Name: "default", Password: generateUUID()},
	// 		}, "grpc", "trojan"),
	// 		// BuildInboundConfig("socks","socks-in",8093,tcpFastOpen,domainStrategy,listen,sniff,sniffTimeout,[]User{
	// 		// 	{Username: "default",Password: generateUUID()},
	// 		// },"",""),
	// 		// {
	// 		// 	Type: "socks",
	// 		// 	Tag: "socks-in",
	// 		// 	Listen: "0.0.0.0",
	// 		// 	ListenPort: 8093,
	// 		// 	TCPFastOpen: true,
	// 		// 	DomainStrategy: domainStrategy,
	// 		// 	Sniff: true,
	// 		// 	SniffTimeout: sniffTimeout,
	// 		// 	SniffOverrideDestination: false,
	// 		// 	Users: []User{
	// 		// 		{Username: "default", Password: generateUUID()},
	// 		// 	},
	// 		// 	Multiplex: Multiplex{},
	// 		// 	Transport: Transport{},
	// 		// },
	// 	},
	// 	Outbounds: createOutbounds(),
	// }

	logDNSExperimental := createLogDNSExperimental()

	if err := logDNSExperimental.SaveToFile(filepath.Join("/config.json")); err != nil {
		log.Fatalf("failed to save config: %v", err)
	}

	log.Print("config.json created successfully.")
}

func createLogDNSExperimental() Config{
	return Config{
		Log: Log{
			Level:     "info",
			Output:    "/usr/local/sing-box/sing-box.log",
			Timestamp: true,
		},
		NTP: NTP{
			Interval:   "5m0s",
			Server:     "time.apple.com",
			ServerPort: 123,
			Detour:     "direct",
		},
		Experimental: Experimental{
			CacheFile: CacheFile{
				Enabled:    true,
				Path:       "cache.db",
				StoreRDrc: true,
				CacheID:    "sing-box",
			},
			ClashAPI: ClashAPI{
				ExternalController:        "0.0.0.0:9090",
				ExternalUI:                "yacd",
				ExternalUIDownloadURL:     "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip",
				ExternalUIDownloadDetour:  "direct",
				Secret: "qwe12345",
			},
		},
		DNS: BuildDNSConfig(),
	}
}

func createInbound() Config{
	listen := "0.0.0.0"
	domainStrategy := "prefer_ipv4"
	tcpFastOpen := true
	sniffTimeout := "300ms"
	sniff := true
	return Config{
		Inbounds: []Inbound{
			BuildInboundConfig("vless", "vless-ws-in", 8001, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "ws", "vless"),
			BuildInboundConfig("vmess", "vmess-ws-in", 8002, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "ws", "vmess"),
			BuildInboundConfig("trojan", "trojan-ws-in", 8003, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", Password: generateUUID()},
			}, "ws", "trojan"),
			BuildInboundConfig("vless", "vless-httpupgrade-in", 8004, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "httpupgrade", "vless"),
			BuildInboundConfig("vmess", "vmess-httpupgrade-in", 8005, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "httpupgrade", "vmess"),
			BuildInboundConfig("trojan", "trojan-httpupgrade-in", 8006, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", Password: generateUUID()},
			}, "httpupgrade", "trojan"),
			BuildInboundConfig("vless", "vless-grpc-in", 8007, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "grpc", "vless"),
			BuildInboundConfig("vmess", "vmess-grpc-in", 8008, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", UUID: generateUUID()},
			}, "grpc", "vmess"),
			BuildInboundConfig("trojan", "trojan-grpc-in", 8009, tcpFastOpen, domainStrategy, listen, sniff, sniffTimeout, []User{
				{Name: "default", Password: generateUUID()},
			}, "grpc", "trojan"),
			BuildInboundConfig("socks","socks-in",8093,tcpFastOpen,domainStrategy,listen,sniff,sniffTimeout,[]User{
				{Name: "default",Password: generateUUID()},
			},"",""),
		},
	}
}

func createOutbounds() Config{
	return Config{
		Outbounds: []Outbound{
			BuildOutboundConfig("direct", "direct", nil, ""),
			BuildOutboundConfig("block", "block", nil, ""),
			BuildOutboundConfig("dns-out", "dns", nil, ""),
			BuildOutboundConfig("TrafficUmum", "selector", []string{"direct"}, "direct"),
			BuildOutboundConfig("TrafficGoogle", "selector", []string{"direct"}, "direct"),
			BuildOutboundConfig("TrafficAds", "selector", []string{"direct", "block"}, "block"),
			BuildOutboundConfig("TrafficPorn", "selector", []string{"direct", "block"}, "block"),
		},
	}
}