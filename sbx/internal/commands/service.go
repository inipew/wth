package commands

import (
	"fmt"
	cli "sbx/internal/cmd"
	"sbx/internal/logger"
	"sbx/pkg/systemdservice"
)

func CreateServiceCommand() *cli.Command {
	return &cli.Command{
		Name:        "service",
		Description: "Create systemd services",
		Run: func(cmd *cli.Command, args []string) error {
			return createServices()
		},
		Help: "The 'service' command generates systemd services for Caddy and sing-box.",
	}
}

func createServices() error {
	serviceManager, err := systemdservice.NewServiceManager()
	if err != nil {
		return fmt.Errorf("failed to create service manager: %w", err)
	}

	if err := serviceManager.GenerateCaddyService(); err != nil {
		return fmt.Errorf("failed to generate Caddy service: %w", err)
	}

	if err := serviceManager.GenerateSingBoxService(); err != nil {
		return fmt.Errorf("failed to generate SingBox service: %w", err)
	}

	logger.GetLogger().Info().Msg("Services generated successfully")
	return nil
}