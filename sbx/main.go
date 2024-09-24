package main

import (
	"context"
	"flag"
	"path/filepath"
	"sbx/internal/archive"
	"sbx/internal/caddyfile"
	"sbx/internal/cmd"
	"sbx/internal/config"
	"sbx/internal/download"
	"sbx/internal/github"
	"sbx/internal/logger"
	"sbx/internal/systemdservice"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	githubTimeout    = 10 * time.Second
	apiTimeout       = 30 * time.Second
	downloadTimeout  = 1 * time.Hour
)

func main() {
	if err := logger.InitLogger(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	myCLI := cmd.NewCLI("1.0.0")

	// Register commands
	myCLI.AddCommand(createSingCommand())
	myCLI.AddCommand(createCaddyCommand())
	myCLI.AddCommand(createCaddyfileCommand())
	myCLI.AddCommand(createServiceCommand())

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
			if err := performDownloadSing(*latestFlag); err != nil {
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
			if err := performDownloadCaddy(*latestFlag); err != nil {
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

// performDownloadCaddy handles the download of Caddy
func performDownloadCaddy(preRelease bool) error {
	return performDownload("caddyserver", "caddy", preRelease)
}

// performDownloadSing handles the download of sing-box
func performDownloadSing(preRelease bool) error {
	return performDownload("SagerNet", "sing-box", preRelease)
}

// performDownload downloads a repository
func performDownload(repoOwner, repoName string, preRelease bool) error {
	client := github.NewClient(githubTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	filePath := filepath.Join(config.TmpDir, repoName+".tar.gz")

	version, err := client.GetLatestRelease(ctx, repoOwner, repoName, preRelease)
	if err != nil {
		return err
	}

	downloadURL, err := github.BuildDownloadURL(repoOwner, repoName, version)
	if err != nil {
		return err
	}

	log.Info().Str("version", version).Str("url", downloadURL).Msg("Download information")

	downloadClient := download.NewClient(
		30*time.Second, // timeout
		5,              // retryCount
		5*time.Second,  // retryDelay
		10,             // concurrentChunks
		5*1024*1024,    // chunkSize (5MB)
	)

	ctx, cancel = context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	if err := downloadClient.DownloadFile(ctx, downloadURL, filePath); err != nil {
		return err
	}

	if err := archive.UntarGz(filePath, filepath.Dir(config.CaddyBinPath)); err != nil {
		return err
	}

	return nil
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
