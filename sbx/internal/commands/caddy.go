package commands

import (
	"flag"
	cli "sbx/internal/cmd"
	"sbx/pkg/downloader"
)

func CreateCaddyCommand() *cli.Command {
	flags := &cli.FlagSetParser{FlagSet: flag.NewFlagSet("caddy", flag.ContinueOnError)}
	latestFlag := flags.Bool("latest", false, "Download the latest prerelease version")

	return &cli.Command{
		Name:        "caddy",
		Description: "Download the latest release of Caddy",
		Flags:       flags,
		Run: func(cmd *cli.Command, args []string) error {
			return downloader.PerformDownloadCaddy(*latestFlag)
		},
		Help:    "The 'caddy' command downloads the latest release of Caddy. Use '--latest' for the prerelease version.",
		Aliases: []string{"c"},
	}
}