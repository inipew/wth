package account

import (
	"fmt"
	"singconfig/internal/config/inbound"
	"singconfig/pkg/singbox"
	"strings"
)

type InboundGroups struct {
    Remarks     string
    Ports       Ports
    Username    string
    UUID        string
    Path        string
    ServiceName string
    Networks    map[string]bool
    Type        string
}

type Ports struct {
    PortsTLS    string
    PortsNTLS   string
    Port        string
}

var protocolOrder = []string{"vless", "vmess", "trojan", "socks"}

func PrintInboundConfigs(config *singbox.SingBoxConfig) {
    if config == nil {
        fmt.Println("Error: SingBoxConfig is nil")
        return
    }

    if len(config.Inbounds) == 0 {
        fmt.Println("No inbounds found in the configuration")
        return
    }

    groupedConfig := make(map[string]*InboundGroups)

    for i, inbound := range config.Inbounds {
        if len(inbound.Users) == 0 {
            fmt.Printf("No users found for inbound %d\n", i)
            continue
        }

        for _, user := range inbound.Users {
            key := generateKeys(user, inbound.Type)
            group, exists := groupedConfig[key]
            if !exists {
                group = createGroup(inbound, user)
                groupedConfig[key] = group
            }
            updateGroup(group, inbound)
        }
    }

    if len(groupedConfig) == 0 {
        fmt.Println("No groups were created. Check if the input data is correct.")
        return
    }

    printGroupedConfigs(groupedConfig)
}

func generateKeys(user inbound.UserConfig, inboundType string) string {
    if user.UUID != "" {
        return fmt.Sprintf("%s_%s", user.UUID, inboundType)
    }
    return fmt.Sprintf("%s_%s", user.Password, inboundType)
}

func createGroup(inbound inbound.InboundConfig, user inbound.UserConfig) *InboundGroups {
    name := user.Name
    if name == "" {
        name = user.Username
    }

	pass := user.UUID
    if pass == "" {
        pass = user.Password
    }
    var port string
    if inbound.Type == "socks" {
        port = fmt.Sprintf("%d", inbound.ListenPort)
    }

    return &InboundGroups{
        Remarks:  fmt.Sprintf("%s_%s", inbound.Type, name),
        Ports: Ports{
            PortsTLS: "443,2053,2083,2087,2096,8443",
            PortsNTLS: "80,8080,8880,2052,2082,2086,2095",
            Port: port,
        },
        Username: name,
        UUID:     pass,
        Networks: make(map[string]bool),
        Type:     inbound.Type,  // Simpan tipe protokol
    }
}

func updateGroup(group *InboundGroups, inbound inbound.InboundConfig) {
    if inbound.Transport != nil {
        if inbound.Transport.Path != "" {
            group.Path = inbound.Transport.Path
        }
        if inbound.Transport.ServiceName != "" {
            group.ServiceName = inbound.Transport.ServiceName
        }
        if inbound.Transport.Type != "" {
            group.Networks[inbound.Transport.Type] = true
        }
    }
}

func printGroupedConfigs(groupedConfig map[string]*InboundGroups) {
    for _, protocol := range protocolOrder {
        printed := false
        for _, group := range groupedConfig {
            if group.Type == protocol {
                if !printed {
                    fmt.Printf("=== %s ===\n", strings.ToUpper(protocol))
                    printed = true
                }
                printGroup(group)
                fmt.Println()
            }
        }
    }
}

func printGroup(group *InboundGroups) {
    var b strings.Builder

    b.WriteString(fmt.Sprintf("Remarks\t\t: %s\n", group.Remarks))
    if group.Type == "socks" {
        b.WriteString(fmt.Sprintf("Port TLS\t: %s\n", group.Ports.Port))
    } else {
        b.WriteString(fmt.Sprintf("Port TLS\t: %s\n", group.Ports.PortsTLS))
        b.WriteString(fmt.Sprintf("Port nTLS\t: %s\n", group.Ports.PortsNTLS))
    }
    if group.Username != "" {
        b.WriteString(fmt.Sprintf("Username\t: %s\n", group.Username))
    }
    if group.UUID != "" {
        b.WriteString(fmt.Sprintf("UUID\t\t: %s\n", group.UUID))
    }
    if len(group.Networks) > 0 {
        networks := make([]string, 0, len(group.Networks))
        for network := range group.Networks {
            networks = append(networks, network)
        }
        b.WriteString(fmt.Sprintf("Network\t\t: %s\n", strings.Join(networks, ", ")))
    }
    if group.Path != "" {
        b.WriteString(fmt.Sprintf("Path\t\t: %s\n", group.Path))
    }
    if group.ServiceName != "" {
        b.WriteString(fmt.Sprintf("ServiceName\t: %s\n", group.ServiceName))
    }

    fmt.Print(b.String())
}