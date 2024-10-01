package singbox

import (
	"encoding/json"
	"os"
	"singconfig/internal/config/dns"
	"singconfig/internal/config/experimental"
	"singconfig/internal/config/inbound"
	"singconfig/internal/config/log"
	"singconfig/internal/config/ntp"
	"singconfig/internal/config/outbound"
	"singconfig/internal/config/route"
)

type SingBoxConfig struct {
	Log          *log.LogConfig						`json:"log,omitempty"`
	DNS          *dns.DNSConfig						`json:"dns,omitempty"`
	NTP          *ntp.NTPConfig						`json:"ntp,omitempty"`
	Experimental *experimental.ExperimentalConfig	`json:"experimental,omitempty"`
	Inbounds     []inbound.InboundConfig			`json:"inbounds,omitempty"`
	Outbounds    []outbound.OutboundConfig			`json:"outbounds,omitempty"`
	Route        *route.RouteConfig					`json:"route,omitempty"`
}

// func (c *SingBoxConfig) SaveToFile(filename string) error {
// 	data, err := json.MarshalIndent(c, "", "  ")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(filename, data, 0644)
// }
func (c *SingBoxConfig) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}