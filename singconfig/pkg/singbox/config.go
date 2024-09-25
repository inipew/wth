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
	Log          *log.LogConfig
	DNS          *dns.DNSConfig
	NTP          *ntp.NTPConfig
	Experimental *experimental.ExperimentalConfig
	Inbounds     []inbound.InboundConfig
	Outbounds    []outbound.OutboundConfig
	Route        *route.RouteConfig
}

func (c *SingBoxConfig) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// File: pkg/singbox/v1/config.go

// package singbox

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"

// 	"singconfig/internal/interfaces"
// )

// // SingBoxConfig represents the main configuration structure for SingBox
// type SingBoxConfig struct {
// 	Log          interface{}   `json:"log,omitempty"`
// 	DNS          interface{}   `json:"dns,omitempty"`
// 	NTP          interface{}   `json:"ntp,omitempty"`
// 	Experimental interface{}   `json:"experimental,omitempty"`
// 	Inbounds     []interface{} `json:"inbounds,omitempty"`
// 	Outbounds    []interface{} `json:"outbounds,omitempty"`
// 	Route        interface{}   `json:"route,omitempty"`
// }

// // NewSingBoxConfig creates a new instance of SingBoxConfig
// func NewSingBoxConfig() *SingBoxConfig {
// 	return &SingBoxConfig{
// 		Inbounds:  make([]interface{}, 0),
// 		Outbounds: make([]interface{}, 0),
// 	}
// }

// // SetLog sets the log configuration
// func (c *SingBoxConfig) SetLog(log interface{}) interfaces.SingBoxConfig {
// 	c.Log = log
// 	return c
// }

// // SetDNS sets the DNS configuration
// func (c *SingBoxConfig) SetDNS(dns interface{}) interfaces.SingBoxConfig {
// 	c.DNS = dns
// 	return c
// }

// // SetNTP sets the NTP configuration
// func (c *SingBoxConfig) SetNTP(ntp interface{}) interfaces.SingBoxConfig {
// 	c.NTP = ntp
// 	return c
// }

// // AddInbound adds an inbound configuration
// func (c *SingBoxConfig) AddInbound(inbound interface{}) interfaces.SingBoxConfig {
// 	c.Inbounds = append(c.Inbounds, inbound)
// 	return c
// }

// // AddOutbound adds an outbound configuration
// func (c *SingBoxConfig) AddOutbound(outbound interface{}) interfaces.SingBoxConfig {
// 	c.Outbounds = append(c.Outbounds, outbound)
// 	return c
// }

// // SetRoute sets the route configuration
// func (c *SingBoxConfig) SetRoute(route interface{}) interfaces.SingBoxConfig {
// 	c.Route = route
// 	return c
// }

// // SetExperimental sets the experimental configuration
// func (c *SingBoxConfig) SetExperimental(experimental interface{}) interfaces.SingBoxConfig {
// 	c.Experimental = experimental
// 	return c
// }

// // SaveToFile saves the configuration to a file
// func (c *SingBoxConfig) SaveToFile(filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %w", err)
// 	}
// 	defer file.Close()

// 	return c.SaveToWriter(file)
// }

// // SaveToWriter saves the configuration to an io.Writer
// func (c *SingBoxConfig) SaveToWriter(writer io.Writer) error {
// 	encoder := json.NewEncoder(writer)
// 	encoder.SetIndent("", "  ")
// 	if err := encoder.Encode(c); err != nil {
// 		return fmt.Errorf("failed to encode config: %w", err)
// 	}
// 	return nil
// }

// // LoadFromFile loads the configuration from a file
// func LoadFromFile(filename string) (*SingBoxConfig, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open file: %w", err)
// 	}
// 	defer file.Close()

// 	return LoadFromReader(file)
// }

// // LoadFromReader loads the configuration from an io.Reader
// func LoadFromReader(reader io.Reader) (*SingBoxConfig, error) {
// 	var config SingBoxConfig
// 	decoder := json.NewDecoder(reader)
// 	if err := decoder.Decode(&config); err != nil {
// 		return nil, fmt.Errorf("failed to decode config: %w", err)
// 	}
// 	return &config, nil
// }