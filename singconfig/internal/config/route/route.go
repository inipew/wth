package route

const (
	defaultUpdateInterval = "168h0m0s"
	cdnBaseURL            = "https://cdn.jsdelivr.net/gh/"
)

type RouteConfig struct {
	Rules               []RouteRuleConfig `json:"rules,omitempty"`
	RuleSet             []RuleSetConfig   `json:"rule_set,omitempty"`
	Final               string            `json:"final,omitempty"`
	AutoDetectInterface bool              `json:"auto_detect_interface,omitempty"`
}

type RouteRuleConfig struct {
	Type         string            `json:"type,omitempty"`
	Mode         string            `json:"mode,omitempty"`
	Rules        []RouteRuleConfig `json:"rules,omitempty"`
	Outbound     string            `json:"outbound,omitempty"`
	Port         int               `json:"port,omitempty"`
	Protocol     string            `json:"protocol,omitempty"`
	RuleSet      []string          `json:"rule_set,omitempty"`
	DomainSuffix []string          `json:"domain_suffix,omitempty"`
}

type RuleSetConfig struct {
	Type            string `json:"type"`
	Tag             string `json:"tag"`
	Format          string `json:"format"`
	URL             string `json:"url"`
	DownloadDetour  string `json:"download_detour"`
	UpdateInterval  string `json:"update_interval"`
}

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
		newRule().withDNS(),
		newRule().withRuleSet([]string{"direct_some_web", "geosite-google_ads"}).withOutbound("direct"),
		newRule().withRuleSet([]string{"geosite-malicious", "geoip-malicious"}).withOutbound("block"),
		newRule().withQuicYoutube(),
		newRule().withRuleSet([]string{"oisd-full", "rule-ads", "d3ward"}).withOutbound("TrafficAds"),
		newRule().withRuleSet([]string{"oisd-nsfw", "category-porn"}).withOutbound("TrafficPorn"),
		newRule().withDomainSuffix([]string{"gstatic.com", "cp.cloudflare.com"}).withOutbound("TrafficUmum"),
		newRule().withRuleSet([]string{"geoip_google", "google", "rule_youtube"}).withOutbound("TrafficGoogle"),
		newRule().withRuleSet([]string{"openai", "needwarp", "reddit", "geoip_cloudflare", "microsoft", "rule_github"}).withOutbound("direct"),
		newRule().withRuleSet([]string{"googlefcm", "commonports"}).withOutbound("TrafficUmum"),
	}
}

func newRule() RouteRuleConfig {
	return RouteRuleConfig{}
}

func (r RouteRuleConfig) withDNS() RouteRuleConfig {
	r.Type = "logical"
	r.Mode = "or"
	r.Rules = []RouteRuleConfig{
		{Port: 53},
		{Protocol: "dns"},
	}
	r.Outbound = "dns-out"
	return r
}

func (r RouteRuleConfig) withRuleSet(ruleSet []string) RouteRuleConfig {
	r.RuleSet = ruleSet
	return r
}

func (r RouteRuleConfig) withOutbound(outbound string) RouteRuleConfig {
	r.Outbound = outbound
	return r
}

func (r RouteRuleConfig) withDomainSuffix(domainSuffix []string) RouteRuleConfig {
	r.DomainSuffix = domainSuffix
	return r
}

func (r RouteRuleConfig) withQuicYoutube() RouteRuleConfig {
	r.Type = "logical"
	r.Mode = "and"
	r.Rules = []RouteRuleConfig{
		{Protocol: "quic"},
		{RuleSet: []string{"rule_youtube"}},
	}
	r.Outbound = "block"
	return r
}

func createRuleSets() []RuleSetConfig {
	ruleSets := []struct {
		tag            string
		url            string
		updateInterval string
	}{
		{"direct_some_web", "inipew/any@main/direct-some-web.srs", "24h0m0s"},
		{"needwarp", "inipew/any@main/warped.srs", "24h0m0s"},
		{"adguard", "inipew/any@main/adguard.srs", "24h0m0s"},
		{"geoip-malicious", "malikshi/sing-box-geo@rule-set-geoip/geoip-malicious.srs", defaultUpdateInterval},
		{"oisd-full", "malikshi/sing-box-geo@rule-set-geosite/geosite-oisd-full.srs", defaultUpdateInterval},
		{"oisd-nsfw", "malikshi/sing-box-geo@rule-set-geosite/geosite-oisd-nsfw.srs", defaultUpdateInterval},
		{"rule-ads", "malikshi/sing-box-geo@rule-set-geosite/geosite-rule-ads.srs", defaultUpdateInterval},
		{"d3ward", "malikshi/sing-box-geo@rule-set-geosite/geosite-d3ward.srs", defaultUpdateInterval},
		{"geosite-malicious", "malikshi/sing-box-geo@rule-set-geosite/geosite-rule-malicious.srs", defaultUpdateInterval},
		{"category-porn", "malikshi/sing-box-geo@rule-set-geosite/geosite-category-porn.srs", defaultUpdateInterval},
		{"openai", "malikshi/sing-box-geo@rule-set-geosite/geosite-openai.srs", defaultUpdateInterval},
		{"onedrive", "malikshi/sing-box-geo@rule-set-geosite/geosite-onedrive.srs", defaultUpdateInterval},
		{"microsoft", "malikshi/sing-box-geo@rule-set-geosite/geosite-microsoft.srs", defaultUpdateInterval},
		{"rule_github", "malikshi/sing-box-geo@rule-set-geosite/geosite-github.srs", defaultUpdateInterval},
		{"rule_youtube", "malikshi/sing-box-geo@rule-set-geosite/geosite-youtube.srs", defaultUpdateInterval},
		{"google", "malikshi/sing-box-geo@rule-set-geosite/geosite-google.srs", defaultUpdateInterval},
		{"googlefcm", "malikshi/sing-box-geo@rule-set-geosite/geosite-googlefcm.srs", defaultUpdateInterval},
		{"geoip_google", "malikshi/sing-box-geo@rule-set-geoip/geoip-google.srs", defaultUpdateInterval},
		{"geosite-google_ads", "malikshi/sing-box-geo@rule-set-geosite/geosite-google-ads.srs", defaultUpdateInterval},
		{"reddit", "malikshi/sing-box-geo@rule-set-geosite/geosite-reddit.srs", defaultUpdateInterval},
		{"geoip_cloudflare", "malikshi/sing-box-geo@rule-set-geoip/geoip-cloudflare.srs", defaultUpdateInterval},
		{"commonports", "inipew/any@main/commonports.srs", "720h0m0s"},
	}

	var result []RuleSetConfig
	for _, rs := range ruleSets {
		result = append(result, createRemoteRuleSet(rs.tag, cdnBaseURL+rs.url, rs.updateInterval))
	}
	return result
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