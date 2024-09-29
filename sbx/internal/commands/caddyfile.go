package commands

import (
	"fmt"
	"sbx/internal/caddyfile"
	cli "sbx/internal/cmd"
	"sbx/internal/config"
	"sbx/internal/logger"
	"sbx/internal/templates"
)

func CreateCaddyfileCommand() *cli.Command {
	return &cli.Command{
		Name:        "caddyfile",
		Description: "Generate a Caddyfile",
		Run: func(cmd *cli.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("domain is required for Caddyfile generation")
			}
			return generateCaddyfile(args[0])
		},
		Help: "The 'caddyfile' command generates a Caddyfile for the specified domain. Usage: sbx caddyfile <domain>",
	}
}

func generateCaddyfile(domain string) error {
	caddy := caddyfile.NewCaddyfileGenerator(templates.CaddyfileTemplate, caddyfile.Config{Domain: domain})
	caddyContent, err := caddy.Generate()
	if err != nil {
		return err
	}

	if err := caddyfile.WriteToFile(caddyContent, config.CaddyFilePath()); err != nil {
		return err
	}

	logger.GetLogger().Info().Str("filepath", config.CaddyFilePath()).Msg("Caddyfile generated successfully")
	return nil
}