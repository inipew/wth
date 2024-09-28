package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"singconfig/internal/config/dns"
	"singconfig/internal/config/experimental"
	"singconfig/internal/config/inbound"
	"singconfig/internal/config/log"
	"singconfig/internal/config/ntp"
	"singconfig/internal/config/outbound"
	"singconfig/internal/config/route"
	"singconfig/internal/utils"
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

func SaveConfig(config *singbox.SingBoxConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// AddUser adds a new user to the SingBox configuration based on specified criteria.
func AddUser(config *singbox.SingBoxConfig, transportType, userType, name string) error {
    if config == nil || config.Inbounds == nil {
        return errors.New("config or Inbounds is nil")
    }

    for i := range config.Inbounds {
        if shouldAddUser(&config.Inbounds[i], transportType, userType) {
            if err := addUserToInbound(&config.Inbounds[i], name); err != nil {
                return err
            }
        }
    }

    return nil
}

// shouldAddUser determines if a user should be added to the given inbound configuration.
func shouldAddUser(inboundConfig *inbound.InboundConfig, transportType, userType string) bool {
    matchesUserType := userType == "" || userType == "all" || inboundConfig.Type == userType

    // Special handling for socks, which may not have a Transport configuration
    if inboundConfig.Type == "socks" {
        return matchesUserType && (transportType == "" || transportType == "all")
    }

    // For other types, check if Transport exists and matches
    if inboundConfig.Transport == nil {
        return false
    }

    matchesTransport := transportType == "" || transportType == "all" || inboundConfig.Transport.Type == transportType

    return matchesTransport && matchesUserType
}

// addUserToInbound adds a new user to the inbound configuration.
func addUserToInbound(inboundConfig *inbound.InboundConfig, name string) error {
    uuid := utils.GenerateUUID()
    newUser, err := createNewUser(inboundConfig.Type, name, uuid)
    if err != nil {
        return err
    }

    inboundConfig.Users = append(inboundConfig.Users, newUser)
    return nil
}

// createNewUser creates a new user based on the inbound type.
func createNewUser(inboundType, name, uuid string) (inbound.UserConfig, error) {
    switch inboundType {
    case "vmess", "vless":
        return inbound.UserConfig{
            Name: name,
            UUID: uuid,
        }, nil
    case "trojan":
        return inbound.UserConfig{
            Name:     name,
            Password: uuid,
        }, nil
    case "socks":
        return inbound.UserConfig{
            Username: name,
            Password: uuid,
        }, nil
    default:
        return inbound.UserConfig{}, errors.New("unsupported inbound type")
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

// DisplayInboundDetails menampilkan detail dari semua inbound dalam konfigurasi
func DisplayInboundDetails(config *singbox.SingBoxConfig) {
	if config == nil || config.Inbounds == nil {
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