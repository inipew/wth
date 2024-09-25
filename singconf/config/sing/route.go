package sing

func BuildRouteConfig() *RouteConfig {
	return &RouteConfig{
		Rules:               buildRouteRules(),
		RuleSet:             createRuleSets(),
		Final:               "TrafficUmum",
		AutoDetectInterface: true,
	}
}

func buildRouteRules() []RouteRuleConfig {
	return []RouteRuleConfig{
		buildDNSRule(),
		buildDirectRule(),
		buildMaliciousRule(),
		buildQuicYoutubeRule(),
		buildAdsRule(),
		buildPornRule(),
		buildCloudflareRule(),
		buildGoogleRule(),
		buildWarpRule(),
		buildCommonPortsRule(),
	}
}

func buildDNSRule() RouteRuleConfig {
	return RouteRuleConfig{
		Type: "logical",
		Mode: "or",
		Rules: []RouteRuleConfig{
			{Port: 53},
			{Protocol: "dns"},
		},
		Outbound: "dns-out",
	}
}

func buildDirectRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"direct_some_web", "geosite-google_ads"},
		Outbound: "direct",
	}
}

func buildMaliciousRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"geosite-malicious", "geoip-malicious"},
		Outbound: "block",
	}
}

func buildQuicYoutubeRule() RouteRuleConfig {
	return RouteRuleConfig{
		Type: "logical",
		Mode: "and",
		Rules: []RouteRuleConfig{
			{Protocol: "quic"},
			{RuleSet: []string{"rule_youtube"}},
		},
		Outbound: "block",
	}
}

func buildAdsRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"oisd-full", "rule-ads", "d3ward"},
		Outbound: "TrafficAds",
	}
}

func buildPornRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"oisd-nsfw", "category-porn"},
		Outbound: "TrafficPorn",
	}
}

func buildCloudflareRule() RouteRuleConfig {
	return RouteRuleConfig{
		DomainSuffix: []string{"gstatic.com", "cp.cloudflare.com"},
		Outbound:     "TrafficUmum",
	}
}

func buildGoogleRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"geoip_google", "google", "rule_youtube"},
		Outbound: "TrafficGoogle",
	}
}

func buildWarpRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"openai", "needwarp", "reddit", "geoip_cloudflare", "microsoft", "rule_github"},
		Outbound: "direct",
	}
}

func buildCommonPortsRule() RouteRuleConfig {
	return RouteRuleConfig{
		RuleSet:  []string{"googlefcm", "commonports"},
		Outbound: "TrafficUmum",
	}
}

func createRuleSets() []RuleSetConfig {
	return []RuleSetConfig{
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

func createRemoteRuleSet(tag, url, updateInterval string) RuleSetConfig {
	return RuleSetConfig{
		Type:           "remote",
		Tag:            tag,
		Format:         "binary",
		URL:            url,
		DownloadDetour: "direct",
		UpdateInterval: updateInterval,
	}
}