package account

import (
	"errors"
	"singconfig/internal/config/inbound"
	"singconfig/internal/utils"
	"singconfig/pkg/singbox"
)

// AddUser adds a new user to the SingBox configuration based on specified criteria.
func AddUser(config *singbox.SingBoxConfig, transportType, userType, name string, uuid string) error {
	if config == nil || config.Inbounds == nil {
		return errors.New("config or Inbounds is nil")
	}
	if name == "" || uuid == "" {
		return errors.New("name and uuid cannot be empty")
	}
	if uuid == "" {
		uuid = utils.GenerateUUID()
	}
	for i := range config.Inbounds {
		if shouldAddUserToInbound(&config.Inbounds[i], transportType, userType) {
			if err := addUserToInbound(&config.Inbounds[i], name, uuid); err != nil {
				return err
			}
		}
	}
	return nil
}

// shouldAddUserToInbound determines if a user should be added to the given inbound configuration.
func shouldAddUserToInbound(inboundConfig *inbound.InboundConfig, transportType, userType string) bool {
	matchesUserType := userType == "" || userType == "all" || inboundConfig.Type == userType
	if inboundConfig.Type == "socks" {
		return matchesUserType && (transportType == "" || transportType == "all")
	}
	if inboundConfig.Transport == nil {
		return false
	}
	matchesTransport := transportType == "" || transportType == "all" || inboundConfig.Transport.Type == transportType
	return matchesTransport && matchesUserType
}

// addUserToInbound adds a new user to the inbound configuration.
func addUserToInbound(inboundConfig *inbound.InboundConfig, name string, uuid string) error {
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