package commands

import (
	"flag"
	cli "sbx/internal/cmd"
	"sbx/pkg/downloader"
)

func CreateSingCommand() *cli.Command {
	flags := &cli.FlagSetParser{FlagSet: flag.NewFlagSet("sing", flag.ContinueOnError)}
	latestFlag := flags.Bool("latest", false, "Download the latest prerelease version")

	return &cli.Command{
		Name:        "sing",
		Description: "Download the latest release of sing-box",
		Flags:       flags,
		Run: func(cmd *cli.Command, args []string) error {
			return downloader.PerformDownloadSing(*latestFlag)
		},
		Help:    "The 'sing' command downloads the latest release of sing-box. Use '--latest' for the prerelease version.",
		Aliases: []string{"s"},
	}
}

func setupCommandFlags(commandName string, latestFlag *bool) *cli.FlagSetParser {
	flags := &cli.FlagSetParser{FlagSet: flag.NewFlagSet(commandName, flag.ContinueOnError)}
	flags.BoolVar(latestFlag, "latest", false, "Download the latest prerelease version")
	flags.BoolVar(latestFlag, "l", false, "Download the latest prerelease version (shorthand for --latest")
	return flags
}