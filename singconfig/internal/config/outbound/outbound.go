package outbound

type OutboundConfig struct {
	Type                      string     `json:"type"`
	Tag                       string     `json:"tag"`
	Outbounds                 []string   `json:"outbounds,omitempty"`
	Default                   string     `json:"default,omitempty"`
	InterruptExistConnections bool       `json:"interrupt_exist_connections,omitempty"`
	Detour                    string     `json:"detour,omitempty"`
	DomainStrategy            string     `json:"domain_strategy,omitempty"`
	Interval                  string     `json:"interval,omitempty"`
	IdleTimeout               string     `json:"idle_timeout,omitempty"`
}

func BuildOutbounds() []OutboundConfig {
	outboundTypes := []struct {
		tag       string
		outType   string
		outbounds []string
		defaultOB string
	}{
		{"direct", "direct", nil, ""},
		{"block", "block", nil, ""},
		{"dns-out", "dns", nil, ""},
		{"TrafficUmum", "selector", []string{"direct"}, "direct"},
		{"TrafficGoogle", "selector", []string{"direct"}, "direct"},
		{"TrafficAds", "selector", []string{"direct", "block"}, "block"},
		{"TrafficPorn", "selector", []string{"direct", "block"}, "block"},
	}

	var outbounds []OutboundConfig
	for _, entry := range outboundTypes {
		outbounds = append(outbounds, buildOutboundConfig(entry))
	}
	return outbounds
}

func buildOutboundConfig(entry struct {
	tag       string
	outType   string
	outbounds []string
	defaultOB string
}) OutboundConfig {
	config := OutboundConfig{
		Type: entry.outType,
		Tag:  entry.tag,
	}

	if entry.outType == "selector" {
		config.Outbounds = entry.outbounds
		config.Default = entry.defaultOB
		config.InterruptExistConnections = true
	}

	return config
}