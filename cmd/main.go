package main

import (
	"fmt"
	"os"
	"strings"
)

type Command struct {
    Name        string
    Description string
    Run         func(args []string)
}

type CLI struct {
    commands map[string]Command
}

func NewCLI() *CLI {
    return &CLI{commands: make(map[string]Command)}
}

func (cli *CLI) AddCommand(cmd Command) {
    cli.commands[cmd.Name] = cmd
}

func (cli *CLI) Execute() {
    if len(os.Args) < 2 {
        fmt.Println("Please specify a command.")
        cli.printHelp()
        return
    }

    commandName := os.Args[1]
    command, exists := cli.commands[commandName]

    if !exists {
        fmt.Printf("Unknown command: %s\n", commandName)
        cli.printHelp()
        return
    }

    command.Run(os.Args[2:])
}

func (cli *CLI) printHelp() {
    fmt.Println("Available commands:")
    for _, cmd := range cli.commands {
        fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
    }
}

func main() {
    cli := NewCLI()

    cli.AddCommand(Command{
        Name:        "hello",
        Description: "Prints hello message",
        Run: func(args []string) {
            fmt.Println("Hello, World!")
        },
    })

    cli.AddCommand(Command{
        Name:        "goodbye",
        Description: "Prints goodbye message",
        Run: func(args []string) {
            fmt.Println("Goodbye, World!")
        },
    })

    cli.AddCommand(Command{
        Name:        "echo",
        Description: "Echoes the input",
        Run: func(args []string) {
            fmt.Println(strings.Join(args, " "))
        },
    })

    cli.Execute()
}
