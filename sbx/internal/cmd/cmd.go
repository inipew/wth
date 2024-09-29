package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// CommandFlags interface to abstract flag parsing behavior.
type CommandFlags interface {
	Parse(args []string) error
	Args() []string
	VisitAll(fn func(*flag.Flag))
}

// Command represents a CLI command with associated flags and a runner function.
type Command struct {
	Name        string
	Description string
	Flags       CommandFlags
	Run         func(cmd *Command, args []string) error
	Help        string
	Aliases     []string
}

// CLI represents the command-line interface with registered commands and version information.
type CLI struct {
	commands    map[string]*Command
	version     string
	description string
	rootCommand *Command
}

// NewCLI initializes a new CLI instance with a specified version and description.
func NewCLI(version, description string) *CLI {
	return &CLI{
		commands:    make(map[string]*Command),
		version:     version,
		description: description,
	}
}

// AddCommand registers a command with the CLI, including any aliases.
func (cli *CLI) AddCommand(cmd *Command) {
	cli.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		cli.commands[alias] = cmd
	}
}

// SetRootCommand sets a default command to run when no subcommand is specified.
func (cli *CLI) SetRootCommand(cmd *Command) {
	cli.rootCommand = cmd
}

// Execute processes the command-line arguments and executes the corresponding command.
func (cli *CLI) Execute() {
	if len(os.Args) < 2 {
		if cli.rootCommand != nil {
			if err := cli.runCommand(cli.rootCommand, os.Args[1:]); err != nil {
				cli.displayError(err.Error())
			}
		} else {
			cli.displayError("Please specify a command.")
			cli.PrintHelp()
		}
		return
	}

	commandName := os.Args[1]
	if commandName == "help" {
		if len(os.Args) > 2 {
			cli.PrintCommandHelp(os.Args[2])
		} else {
			cli.PrintHelp()
		}
		return
	}

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
	flagSet.Usage = func() { cli.PrintCommandHelp(command.Name) }
	if command.Flags != nil {
		command.Flags.VisitAll(func(f *flag.Flag) {
			flagSet.Var(f.Value, f.Name, f.Usage)
		})
	}

	if err := flagSet.Parse(args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	return command.Run(command, flagSet.Args())
}

// PrintHelp displays help information for the CLI.
func (cli *CLI) PrintHelp() {
	fmt.Printf("Version: %s\n", cli.version)
	fmt.Printf("Description: %s\n\n", cli.description)
	fmt.Println("Available commands:")
	
	// Collect and sort unique command names
	var names []string
	uniqueCommands := make(map[string]*Command)
	for _, cmd := range cli.commands {
		if _, exists := uniqueCommands[cmd.Name]; !exists {
			uniqueCommands[cmd.Name] = cmd
			names = append(names, cmd.Name)
		}
	}
	sort.Strings(names)

	// Print commands
	for _, name := range names {
		cmd := uniqueCommands[name]
		aliases := getUniqueAliases(cmd.Name, cmd.Aliases)
		aliasStr := ""
		if len(aliases) > 0 {
			aliasStr = fmt.Sprintf(" (aliases: %s)", strings.Join(aliases, ", "))
		}
		fmt.Printf("  %-15s%s%s\n", name, aliasStr, cmd.Description)
	}
	fmt.Println("\nUse 'help <command>' for more information about a command.")
}

// getUniqueAliases returns a slice of unique aliases, excluding the command name itself
func getUniqueAliases(name string, aliases []string) []string {
	uniqueAliases := make(map[string]bool)
	for _, alias := range aliases {
		if alias != name {
			uniqueAliases[alias] = true
		}
	}
	
	result := make([]string, 0, len(uniqueAliases))
	for alias := range uniqueAliases {
		result = append(result, alias)
	}
	sort.Strings(result)
	return result
}

// PrintCommandHelp displays help information for a specific command.
func (cli *CLI) PrintCommandHelp(commandName string) {
	command, exists := cli.commands[commandName]
	if !exists {
		cli.displayError(fmt.Sprintf("Unknown command: %s", commandName))
		return
	}

	fmt.Printf("Usage: %s [flags] [arguments]\n", commandName)
	fmt.Printf("Description: %s\n", command.Description)
	if command.Help != "" {
		fmt.Println("\nHelp:")
		fmt.Println(command.Help)
	}
	if len(command.Aliases) > 0 {
		fmt.Printf("\nAliases: %s\n", strings.Join(command.Aliases, ", "))
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

// FlagSetParser wraps the standard flag.FlagSet to implement the CommandFlags interface.
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