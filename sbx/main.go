package main

import (
	"flag"
	"sbx/commands"
	"sbx/internal/caddyfile"
	"sbx/internal/cmd"
	"sbx/internal/config"
	"sbx/internal/logger"
	"sbx/pkg/singbox"
	"sbx/pkg/systemdservice"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	config.CreateAllDirs()

	myCLI := cmd.NewCLI("1.0.0")

	// Register commands
	myCLI.AddCommand(createSingCommand())
	myCLI.AddCommand(createCaddyCommand())
	myCLI.AddCommand(createCaddyfileCommand())
	myCLI.AddCommand(createServiceCommand())
	myCLI.AddCommand(createSingConfigCommand())

	// Execute the CLI
	myCLI.Execute()
}

// createSingCommand sets up the 'sing' command
func createSingCommand() *cmd.Command {
	latestFlag := new(bool)
	flags := setupCommandFlags("sing", latestFlag)

	return &cmd.Command{
		Name:        "sing",
		Description: "Download the latest release of sing-box.",
		Flags:       flags,
		Run: func(c *cmd.Command, args []string) {
			if err := commands.PerformDownloadSing(*latestFlag); err != nil {
				log.Fatal().Err(err).Msg("Error in download operation")
			}
		},
		Help: "The `sing` command downloads the latest release of sing-box. Use `--latest` for the prerelease version.",
	}
}

// createCaddyCommand sets up the 'caddy' command
func createCaddyCommand() *cmd.Command {
	latestFlag := new(bool)
	flags := setupCommandFlags("caddy", latestFlag)

	return &cmd.Command{
		Name:        "caddy",
		Description: "Download the latest release of Caddy.",
		Flags:       flags,
		Run: func(c *cmd.Command, args []string) {
			if err := commands.PerformDownloadCaddy(*latestFlag); err != nil {
				log.Fatal().Err(err).Msg("Error in download operation")
			}
		},
		Help: "The `caddy` command downloads the latest release of Caddy. Use `--latest` for the prerelease version.",
	}
}

// createCaddyfileCommand sets up the 'caddyfile' command
func createCaddyfileCommand() *cmd.Command {
	return &cmd.Command{
		Name:        "caddyfile",
		Description: "Generate a Caddyfile.",
		Flags:       &cmd.FlagSetParser{FlagSet: flag.NewFlagSet("caddyfile", flag.ContinueOnError)},
		Run: func(c *cmd.Command, args []string) {
			if len(args) < 1 {
				log.Fatal().Msg("Domain is required for Caddyfile generation")
			}
			if err := generateCaddyfile(args[0]); err != nil {
				log.Fatal().Err(err).Msg("Error in Caddyfile generation")
			}
		},
		Help: "The `caddyfile` command generates a Caddyfile for the specified domain.",
	}
}

func createSingConfigCommand() *cmd.Command {
	return &cmd.Command{
		Name:        "config",
		Description: "Generate sing-box config.",
		Flags:       &cmd.FlagSetParser{FlagSet: flag.NewFlagSet("config", flag.ContinueOnError)},
		Run: func(c *cmd.Command, args []string) {
			
				singbox.CreateSingBoxConfig()
		},
		Help: "The `config` command generates Sing-box Config.",
	}
}

// createServiceCommand sets up the 'service' command
func createServiceCommand() *cmd.Command {
	return &cmd.Command{
		Name:        "service",
		Description: "Create systemd services.",
		Run: func(c *cmd.Command, args []string) {
			if err := createServices(); err != nil {
				log.Fatal().Err(err).Msg("Error in service creation")
			}
		},
		Help: "The `service` command generates systemd services for Caddy and sing-box.",
	}
}

// setupCommandFlags creates common flags for commands
func setupCommandFlags(commandName string, latestFlag *bool) *cmd.FlagSetParser {
	flags := &cmd.FlagSetParser{FlagSet: flag.NewFlagSet(commandName, flag.ContinueOnError)}
	flags.BoolVar(latestFlag, "latest", false, "Download the latest prerelease version")
	flags.BoolVar(latestFlag, "l", false, "Download the latest prerelease version (shorthand for --latest")
	return flags
}

// generateCaddyfile generates a Caddyfile for the given domain
func generateCaddyfile(domain string) error {
	caddyContent, err := caddyfile.GenerateCaddyfile(caddyfile.Config{Domain: domain})
	if err != nil {
		return err
	}

	if err := caddyfile.WriteToFile(caddyContent, config.CaddyFilePath); err != nil {
		return err
	}

	log.Info().Str("filepath", config.CaddyFilePath).Msg("Caddyfile generated successfully")
	return nil
}

// createServices generates systemd services
func createServices() error {
	if err := systemdservice.GenerateCaddyService(); err != nil {
		return err
	}
	if err := systemdservice.GenerateSingBoxService(); err != nil {
		return err
	}

	log.Info().Msg("Services generated successfully")
	return nil
}

// func createDNSJson(){
// 	cfg := singconf.Config{
// 		Log: singconf.Log{
// 			Level:     "info",
// 			Output:    config.SingboxLogFilePath,
// 			Timestamp: true,
// 		},
// 		DNS: singconf.DNSConfig{
// 			Servers: []singconf.DNSServer{
// 				{
// 					Tag:             "remote_dns",
// 					Address:         "https://cloudflare-dns.com/dns-query",
// 					AddressResolver: "dns_local",
// 					Strategy:        "prefer_ipv4",
// 					Detour:          "direct",
// 				},
// 				{
// 					Tag:     "dns_local",
// 					Address: "local",
// 					Strategy: "prefer_ipv4",
// 					Detour:   "direct",
// 				},
// 				{
// 					Tag:    "dns_block",
// 					Address: "rcode://success",
// 				},
// 			},
// 			Rules: []singconf.DNSRule{
// 				{
// 					RuleSet:      []string{"geosite-malicious", "geoip-malicious"},
// 					Server:       "dns_block",
// 					DisableCache: true,
// 				},
// 				{
// 					Type:         "logical",
// 					Mode:         "and",
// 					Rules:       []singconf.Rule{
// 						{Protocol: "quic"},
// 						{RuleSet: "youtube"},
// 					},
// 					Server:       "dns_block",
// 					DisableCache: true,
// 					RewriteTTL:   10,
// 				},
// 				{
// 					Outbounds:     []string{"any"},
// 					Server:       "remote_dns",
// 					ClientSubnet: "103.3.60.0/22",
// 				},
// 			},
// 			Final:            "remote_dns",
// 			IndependentCache: true,
// 		},
// 		NTP: singconf.NTP{
// 			Interval:   "5m0s",
// 			Server:     "time.apple.com",
// 			ServerPort: 123,
// 			Detour:     "direct",
// 		},
// 	}

// 	if err := cfg.SaveToFile(filepath.Join(config.SingboxConfDir,"/00_log_and_dns.json")); err != nil {
// 		log.Fatal().Msgf("failed to save config: %v", err)
// 	}

// 	log.Info().Msg("config.json created successfully.")
// }