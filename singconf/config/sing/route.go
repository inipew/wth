package sing

func NewRouteConfig() *Route {
	return &Route{
		Rules:               createRules(),
		RuleSet:             createRuleSets(),
		Final:               "TrafficUmum",
		AutoDetectInterface: true,
	}
}

func createRules() []RouteRule {
	return []RouteRule{
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

func createDNSRule() RouteRule {
	return RouteRule{
		Type: "logical",
		Mode: "or",
		Rule: []Rule{
			{Port: 53},
			{Protocol: "dns"},
		},
		Outbounds: "dns-out",
	}
}

func createDirectRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"direct_some_web", "geosite-google_ads"},
		Outbounds: "direct",
	}
}

func createMaliciousRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"geosite-malicious", "geoip-malicious"},
		Outbounds: "block",
	}
}

func createQuicYoutubeRule() RouteRule {
	return RouteRule{
		Type: "logical",
		Mode: "and",
		Rule: []Rule{
			{Protocol: "quic"},
			{RuleSet: []string{"rule_youtube"}},
		},
		Outbounds: "block",
	}
}

func createAdsRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"oisd-full", "rule-ads", "d3ward"},
		Outbounds: "TrafficAds",
	}
}

func createPornRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"oisd-nsfw", "category-porn"},
		Outbounds: "TrafficPorn",
	}
}

func createCloudflareRule() RouteRule {
	return RouteRule{
		DomainSuffix: []string{"gstatic.com", "cp.cloudflare.com"},
		Outbounds:    "TrafficUmum",
	}
}

func createGoogleRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"geoip_google", "google", "rule_youtube"},
		Outbounds: "TrafficGoogle",
	}
}

func createWarpRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"openai", "needwarp", "reddit", "geoip_cloudflare", "microsoft", "rule_github"},
		Outbounds: "TrafficWarp",
	}
}

func createCommonPortsRule() RouteRule {
	return RouteRule{
		RuleSet:   []string{"googlefcm", "commonports"},
		Outbounds: "TrafficUmum",
	}
}

func createRuleSets() []RuleSet {
	return []RuleSet{
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

func createRemoteRuleSet(tag, url, updateInterval string) RuleSet {
	return RuleSet{
		Type:           "remote",
		Tag:            tag,
		Format:         "binary",
		URL:            url,
		DownloadDetour: "direct",
		UpdateInterval: updateInterval,
	}
}
// func NewRouteConfig() *Route {
// 	return &Route{
// 			Rules: []RouteRule{
// 				{
// 					Type: "logical",
// 					Mode: "or",
// 					Rule: []Rule{
// 						{Port: 53},
// 						{Protocol: "dns"},
// 					},
// 					Outbounds: "dns-out",
// 				},
// 				{
// 					RuleSet:   []string{"direct_some_web", "geosite-google_ads"},
// 					Outbounds: "direct",
// 				},
// 				{
// 					RuleSet:   []string{"geosite-malicious", "geoip-malicious"},
// 					Outbounds: "block",
// 				},
// 				{
// 					Type: "logical",
// 					Mode: "and",
// 					Rule: []Rule{
// 						{Protocol: "quic"},
// 						{RuleSet: []string{"rule_youtube"}},
// 					},
// 					Outbounds: "block",
// 				},
// 				{
// 					RuleSet:   []string{"oisd-full", "rule-ads", "d3ward"},
// 					Outbounds: "TrafficAds",
// 				},
// 				{
// 					RuleSet:   []string{"oisd-nsfw", "category-porn"},
// 					Outbounds: "TrafficPorn",
// 				},
// 				{
// 					DomainSuffix: []string{"gstatic.com", "cp.cloudflare.com"},
// 					Outbounds:    "TrafficUmum",
// 				},
// 				{
// 					RuleSet:   []string{"geoip_google", "google", "rule_youtube"},
// 					Outbounds: "TrafficGoogle",
// 				},
// 				{
// 					RuleSet:   []string{"openai", "needwarp", "reddit", "geoip_cloudflare", "microsoft", "rule_github"},
// 					Outbounds: "TrafficWarp",
// 				},
// 				{
// 					RuleSet:   []string{"googlefcm", "commonports"},
// 					Outbounds: "TrafficUmum",
// 				},
// 			},
// 			RuleSet: []RuleSet{
// 				createRuleSet("direct_some_web", "24h0m0s"),
// 				createRuleSet("needwarp", "24h0m0s"),
// 				createRuleSet("adguard", "24h0m0s"),
// 				createRuleSet("geosite-google_ads", "24h0m0s"),
// 				createRuleSet("geoip-malicious", "24h0m0s"),
// 				createRuleSet("oisd-full", "168h0m0s"),
// 				createRuleSet("oisd-nsfw", "168h0m0s"),
// 				createRuleSet("rule-ads", "168h0m0s"),
// 				createRuleSet("d3ward", "168h0m0s"),
// 				createRuleSet("geosite-malicious", "168h0m0s"),
// 				createRuleSet("category-porn", "168h0m0s"),
// 				createRuleSet("openai", "168h0m0s"),
// 				createRuleSet("onedrive", "168h0m0s"),
// 				createRuleSet("microsoft", "168h0m0s"),
// 				createRuleSet("rule_github", "168h0m0s"),
// 				createRuleSet("rule_youtube", "168h0m0s"),
// 				createRuleSet("google", "168h0m0s"),
// 				createRuleSet("googlefcm", "168h0m0s"),
// 				createRuleSet("reddit", "168h0m0s"),
// 				createRuleSet("geoip_google", "168h0m0s"),
// 				createRuleSet("geoip_cloudflare", "168h0m0s"),
// 				createRuleSet("commonports", "720h0m0s"),
// 			},
// 			Final:               "TrafficUmum",
// 			AutoDetectInterface: true,
// 		}
// }

// func createRuleSet(tag string, updateInterval string) RuleSet {
// 	return RuleSet{
// 		Type:            "remote",
// 		Tag:             tag,
// 		Format:          "binary",
// 		URL:             "https://cdn.jsdelivr.net/gh/malikshi/sing-box-geo@rule-set-geosite/geosite-" + tag + ".srs",
// 		DownloadDetour:  "direct",
// 		UpdateInterval:  updateInterval,
// 	}
// }