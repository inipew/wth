package config

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
	"singconfig/pkg/singbox"
)

func BuildSingBoxConfig() singbox.SingBoxConfig {
	return singbox.SingBoxConfig{
		Log:          log.BuildLogConfig(),
		DNS:          dns.BuildDNSConfig(),
		NTP:          ntp.BuildNTPConfig(),
		Experimental: experimental.BuildExperimentalConfig(),
		Inbounds:     inbound.BuildInbounds(),
		Outbounds:    outbound.BuildOutbounds(),
		Route:        route.BuildRouteConfig(),
	}
}

func LoadConfig(filename string) (*singbox.SingBoxConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config singbox.SingBoxConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *singbox.SingBoxConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func AddUser(config *singbox.SingBoxConfig, user inbound.UserConfig, addType string, transportType string) {
	for i := range config.Inbounds {
		inbound := &config.Inbounds[i]
		if addType == "all" ||
			(addType == "type" && inbound.Type == user.Name) ||
			(addType == "transport" && inbound.Transport != nil && inbound.Transport.Type == transportType) ||
			(addType == "both" && inbound.Type == user.Name && inbound.Transport != nil && inbound.Transport.Type == transportType) {
			inbound.Users = append(inbound.Users, user)
		}
	}
}

func RemoveUser(config *singbox.SingBoxConfig, username string) {
	for i := range config.Inbounds {
		inbound := &config.Inbounds[i]
		for j := 0; j < len(inbound.Users); j++ {
			if inbound.Users[j].Name == username {
				inbound.Users = append(inbound.Users[:j], inbound.Users[j+1:]...)
				j--
			}
		}
	}
}

func UpdateDNS(config *singbox.SingBoxConfig, address, addressType, strategy string) {
	for i := range config.DNS.Servers {
		server := &config.DNS.Servers[i]
		if server.Tag == "remote_dns" {
			server.Address = formatDNSAddress(address, addressType)
			server.Strategy = strategy
			break
		}
	}
}

func formatDNSAddress(address, addressType string) string {
	switch addressType {
	case "tls":
		return "tls://" + address
	case "https":
		return "https://" + address
	case "h3":
		return "h3://" + address
	default:
		return address
	}
}