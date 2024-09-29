package commands

import cli "sbx/internal/cmd"

func CreateSingConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "config",
		Description: "Generate sing-box config",
		Run: func(cmd *cli.Command, args []string) error {
			return nil
		},
		Help: "The 'config' command generates the Sing-box configuration.",
	}
}