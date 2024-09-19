package main

import (
	"context"
	"flag"
	"path/filepath"
	"sbx/internal/archive"
	"sbx/internal/caddyfile"
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

	// Parse command line flags
	mode := flag.String("mode", "", "Operation mode: sing, caddyfile, or service")
	domain := flag.String("domain", "", "Domain for Caddyfile generation")
	preRelease := flag.Bool("latest", false, "Download the latest prerelease version")
	flag.Parse()

	switch *mode {
	case "sing":
		if err := performDownloadSing(*preRelease); err != nil {
			log.Fatal().Err(err).Msg("Error in download operation")
		}
	case "caddy":
		if err := performDownloadCaddy(*preRelease); err != nil {
			log.Fatal().Err(err).Msg("Error in download operation")
		}
	case "caddyfile":
		if *domain == "" {
			log.Fatal().Msg("Domain is required for Caddyfile generation")
		}
		if err := generateCaddyfile(*domain); err != nil {
			log.Fatal().Err(err).Msg("Error in Caddyfile generation")
		}
	case "service":
		if err := createServices(); err != nil {
			log.Fatal().Err(err).Msg("Error in service creation")
		}
	default:
		log.Fatal().Msg("Invalid or missing mode. Use -mode=sing, -mode=caddy, -mode=caddyfile, or -mode=service")
	}
}

func performDownloadCaddy(preRelease bool) error {
	client := github.NewClient(githubTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	repoOwner := "caddyserver"
	repoName := "caddy"
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
		30*time.Second, // increased timeout
		5,              // increased retryCount
		5*time.Second,  // increased retryDelay
		10,             // increased concurrentChunks
		5*1024*1024,    // increased chunkSize (5MB)
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

func performDownloadSing(preRelease bool) error {
	client := github.NewClient(githubTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	repoOwner := "SagerNet"
	repoName := "sing-box"
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
		30*time.Second, // increased timeout
		5,              // increased retryCount
		5*time.Second,  // increased retryDelay
		10,             // increased concurrentChunks
		5*1024*1024,    // increased chunkSize (5MB)
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

func generateCaddyfile(domain string) error {
	caddyContent, err := caddyfile.GenerateCaddyfile(caddyfile.Config{
		Domain: domain,
	})
	if err != nil {
		return err
	}

	if err := caddyfile.WriteToFile(caddyContent, config.CaddyFilePath); err != nil {
		return err
	}

	log.Info().Str("filepath", config.CaddyFilePath).Msg("Caddyfile generated successfully")
	return nil
}

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
