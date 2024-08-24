package internal

import (
	"flag"
	"fmt"
	"os"
)

type FlagParser interface {
	Parse(args []string) error
	Args() []string
	VisitAll(fn func(*flag.Flag))
}

type Command struct {
	Name        string
	Description string
	Flags       FlagParser
	Run         func(args []string)
	Help        string
	Aliases     []string
}

type CLI struct {
	commands map[string]*Command
	version  string
}

func NewCLI(version string) *CLI {
	return &CLI{commands: make(map[string]*Command), version: version}
}

func (cli *CLI) AddCommand(cmd *Command) {
	cli.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		cli.commands[alias] = cmd
	}
}

func (cli *CLI) Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a command.")
		cli.PrintHelp()
		return
	}

	commandName := os.Args[1]
	command, exists := cli.commands[commandName]

	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		cli.PrintHelp()
		return
	}

	// Create a new flag set for the command
	flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)
	if command.Flags != nil {
		// Copy flags from command.Flags to flagSet
		command.Flags.VisitAll(func(f *flag.Flag) {
			flagSet.Var(f.Value, f.Name, f.Usage)
		})
	}

	// Parse flags
	if len(os.Args) > 2 && os.Args[2] == "--help" {
		cli.PrintCommandHelp(commandName)
		return
	}

	args := os.Args[2:]
	err := flagSet.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing flags: %s\n", err)
		cli.PrintCommandHelp(commandName)
		return
	}

	args = flagSet.Args()
	command.Run(args)
}


func (cli *CLI) PrintHelp() {
	fmt.Printf("CLI Version: %s\n", cli.version)
	fmt.Println("Available commands:")
	for _, cmd := range cli.commands {
		fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("\nUse '<command> --help' for more information about a command.")
}

func (cli *CLI) PrintCommandHelp(commandName string) {
	command, exists := cli.commands[commandName]
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		return
	}

	fmt.Printf("Usage: %s [flags]\n", commandName)
	fmt.Printf("Description: %s\n", command.Description)
	if command.Help != "" {
		fmt.Println("\nHelp:")
		fmt.Println(command.Help)
	}
	fmt.Println("\nFlags:")
	if command.Flags != nil {
		command.Flags.VisitAll(func(f *flag.Flag) {
			fmt.Printf("  --%s: %s (default: %s)\n", f.Name, f.Usage, f.DefValue)
		})
	} else {
		fmt.Println("  No flags available.")
	}
}


func addHelpFlag(flagSet *flag.FlagSet) {
	flagSet.Bool("help", false, "Display help for this command")
}

type FlagSetParser struct {
	*flag.FlagSet
}

func (f *FlagSetParser) Parse(args []string) error {
	return f.FlagSet.Parse(args)
}

func (f *FlagSetParser) Args() []string {
	return f.FlagSet.Args()
}

func (f *FlagSetParser) VisitAll(fn func(*flag.Flag)) {
	f.FlagSet.VisitAll(fn)
}