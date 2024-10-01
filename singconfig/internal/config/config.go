package config

import (
	"encoding/json"
	"fmt"
	"os"
	"singconfig/internal/config/dns"
	"singconfig/internal/config/experimental"
	"singconfig/internal/config/inbound"
	"singconfig/internal/config/log"
	"singconfig/internal/config/ntp"
	"singconfig/internal/config/outbound"
	"singconfig/internal/config/route"
	"singconfig/pkg/singbox"
	"sort"
	"strings"
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

func WriteConfigFile(config *singbox.SingBoxConfig, filename string) error {
    jsonData, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }

    // Gunakan os.WriteFile untuk menulis data ke file
    err = os.WriteFile(filename, jsonData, 0644)
    if err != nil {
        return err
    }

    return nil
}

func WriteJSONConfig(data interface{}, filename string) error {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, jsonData, 0644)
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

func ModifyLogLevel(config *singbox.SingBoxConfig, newLevel string) {
    config.Log.Level = newLevel
}

// DisplayInboundDetails menampilkan detail dari semua inbound dalam konfigurasi
func DisplayInboundDetails(config *singbox.SingBoxConfig) {
	if config == nil || len(config.Inbounds) == 0 {
		fmt.Println("Konfigurasi atau Inbounds tidak tersedia")
		return
	}

	inboundTypes := []string{"vmess", "vless", "trojan", "socks"}

	for _, inboundType := range inboundTypes {
		displayInboundByType(config.Inbounds, inboundType)
	}
}

// displayInboundByType menampilkan detail untuk tipe inbound tertentu
func displayInboundByType(inbounds []inbound.InboundConfig, inboundType string) {
	fmt.Printf("\n%s\n", strings.ToUpper(inboundType))
	
	userMap := make(map[string]map[string]bool)

	for _, inbound := range inbounds {
		if strings.EqualFold(inbound.Type, inboundType) {
			for _, user := range inbound.Users {
				if _, exists := userMap[user.Name]; !exists {
					userMap[user.Name] = make(map[string]bool)
				}
				transportType := getTransportType(inbound.Transport)
				userMap[user.Name][transportType] = true
			}
		}
	}

	fmt.Printf("Jumlah user: %d\n", len(userMap))
	if len(userMap) > 0 {
		fmt.Println("List nama user:")
		// Sortir nama pengguna untuk output yang konsisten
		var sortedNames []string
		for name := range userMap {
			sortedNames = append(sortedNames, name)
		}
		sort.Strings(sortedNames)

		for _, name := range sortedNames {
			protocols := userMap[name]
			var protocolList []string
			for protocol := range protocols {
				protocolList = append(protocolList, protocol)
			}
			sort.Strings(protocolList)
			fmt.Printf("* %s (%s)\n", name, strings.Join(protocolList, ", "))
		}
	} else {
		fmt.Println("Tidak ada user")
	}
}

// getTransportType mengembalikan tipe transport dari konfigurasi transport
func getTransportType(transport *inbound.TransportConfig) string {
	if transport == nil {
		return "default"
	}
	return transport.Type
}