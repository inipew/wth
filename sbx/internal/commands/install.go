package commands

import (
	cli "sbx/internal/cmd"
	"sbx/internal/logger"
	"sbx/pkg/downloader"
)

func CreateInstallCommand() *cli.Command {
	latestFlag := new(bool)
	flags := setupCommandFlags("install", latestFlag)

	return &cli.Command{
		Name:        "install",
		Description: "Install caddy and sing-box.",
		Flags:       flags,
		Run: func(cmd *cli.Command, args []string) error {
			return InstallAll(*latestFlag)
		},
		Help: "The `sing` command downloads the latest release of sing-box. Use `--latest` for the prerelease version.",
	}
}

func InstallAll(latestFlag bool) error{
	if err := downloader.PerformDownloadCaddy(latestFlag); err != nil {
		logger.GetLogger().Fatal().Err(err).Msg("Error in download operation")
	}

	if err := downloader.PerformDownloadSing(latestFlag); err != nil {
		logger.GetLogger().Fatal().Err(err).Msg("Error in download operation")
	}

	if err := createServices(); err != nil {
		logger.GetLogger().Fatal().Err(err).Msg("Error in service operation")
	}
	return nil
}