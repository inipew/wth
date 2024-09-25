package commands

import (
	"context"
	"path/filepath"
	"sbx/internal/archive"
	"sbx/internal/config"
	"sbx/internal/download"
	"sbx/internal/github"
	"sbx/pkg/systemdservice"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	githubTimeout    = 10 * time.Second
	apiTimeout       = 30 * time.Second
	downloadTimeout  = 1 * time.Hour
)

func InstallAll(latestFlag bool) error{
	if err := PerformDownloadCaddy(latestFlag); err != nil {
		log.Fatal().Err(err).Msg("Error in download operation")
	}

	if err := PerformDownloadSing(latestFlag); err != nil {
		log.Fatal().Err(err).Msg("Error in download operation")
	}

	if err := createServices(); err != nil {
		log.Fatal().Err(err).Msg("Error in service creation")
	}
	return nil
}

// performDownloadCaddy handles the download of Caddy
func PerformDownloadCaddy(preRelease bool) error {
	return performDownload("caddyserver", "caddy", preRelease)
}

// performDownloadSing handles the download of sing-box
func PerformDownloadSing(preRelease bool) error {
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