package cli

import (
	"flag"
	"fmt"
	"os"
)

// FlagParser interface to abstract flag parsing behavior.
type FlagParser interface {
	Parse(args []string) error
	Args() []string
	VisitAll(fn func(*flag.Flag))
}

// Command represents a CLI command with associated flags and a runner function.
type Command struct {
	Name        string
	Description string
	Flags       FlagParser
	Run         func(cmd *Command, args []string)
	Help        string
	Aliases     []string
}

// CLI represents the command-line interface with registered commands and version information.
type CLI struct {
	commands map[string]*Command
	version  string
}

// NewCLI initializes a new CLI instance with a specified version.
func NewCLI(version string) *CLI {
	return &CLI{
		commands: make(map[string]*Command),
		version:  version,
	}
}

// AddCommand registers a command with the CLI, including any aliases.
func (cli *CLI) AddCommand(cmd *Command) {
	cli.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		cli.commands[alias] = cmd
	}
}

// Execute processes the command-line arguments and executes the corresponding command.
func (cli *CLI) Execute() {
	if len(os.Args) < 2 {
		cli.displayError("Please specify a command.")
		cli.PrintHelp()
		return
	}

	commandName := os.Args[1]
	command, exists := cli.commands[commandName]

	if !exists {
		cli.displayError(fmt.Sprintf("Unknown command: %s", commandName))
		cli.PrintHelp()
		return
	}

	if err := cli.runCommand(command, os.Args[2:]); err != nil {
		cli.displayError(err.Error())
		cli.PrintCommandHelp(commandName)
	}
}

// runCommand sets up and executes the given command with its flags.
func (cli *CLI) runCommand(command *Command, args []string) error {
	flagSet := flag.NewFlagSet(command.Name, flag.ContinueOnError)
	if command.Flags != nil {
		command.Flags.VisitAll(func(f *flag.Flag) {
			flagSet.Var(f.Value, f.Name, f.Usage)
		})
	}

	// Check for help request
	if len(args) > 0 && args[0] == "--help" {
		cli.PrintCommandHelp(command.Name)
		return nil
	}

	if err := flagSet.Parse(args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	command.Run(command, flagSet.Args())
	return nil
}

// PrintHelp displays help information for the CLI.
func (cli *CLI) PrintHelp() {
	fmt.Printf("CLI Version: %s\n", cli.version)
	fmt.Println("Available commands:")
	for _, cmd := range cli.commands {
		fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("\nUse '<command> --help' for more information about a command.")
}

// PrintCommandHelp displays help information for a specific command.
func (cli *CLI) PrintCommandHelp(commandName string) {
	command, exists := cli.commands[commandName]
	if !exists {
		cli.displayError(fmt.Sprintf("Unknown command: %s", commandName))
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

// displayError prints an error message in a standardized format.
func (cli *CLI) displayError(message string) {
	fmt.Fprintln(os.Stderr, "Error:", message)
}

// FlagSetParser wraps the standard flag.FlagSet to implement the FlagParser interface.
type FlagSetParser struct {
	*flag.FlagSet
}

// Parse processes the provided arguments using the flag set.
func (f *FlagSetParser) Parse(args []string) error {
	return f.FlagSet.Parse(args)
}

// Args returns the non-flag arguments after parsing.
func (f *FlagSetParser) Args() []string {
	return f.FlagSet.Args()
}

// VisitAll iterates over all flags in the flag set.
func (f *FlagSetParser) VisitAll(fn func(*flag.Flag)) {
	f.FlagSet.VisitAll(fn)
}
