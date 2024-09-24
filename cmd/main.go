package main

import (
	"cmd/internal/cli"
	"flag"
	"fmt"
)

func main() {
	cmd := cli.NewCLI("1.0.0")

	// Create the `hello` command
	helloFlags := &cli.FlagSetParser{flag.NewFlagSet("hello", flag.ContinueOnError)}
	helloFlags.String("name", "World", "Name to greet")

	helloCommand := &cli.Command{
		Name:        "hello",
		Description: "Prints a greeting message.",
		Flags:       helloFlags,
		Run: func(c *cli.Command, args []string) {
			name := helloFlags.Lookup("name").Value.String()
			fmt.Printf("Hello, %s!\n", name)
		},
		Help: "The `hello` command prints a greeting message. You can use the `--name` flag to specify the name.",
	}

	cmd.AddCommand(helloCommand)

	// Create the `goodbye` command
	goodbyeFlags := &cli.FlagSetParser{flag.NewFlagSet("goodbye", flag.ContinueOnError)}
	goodbyeFlags.Bool("formal", false, "Use formal goodbye")

	goodbyeCommand := &cli.Command{
		Name:        "goodbye",
		Description: "Prints a farewell message.",
		Flags:       goodbyeFlags,
		Run: func(c *cli.Command, args []string) {
			if goodbyeFlags.Lookup("formal").Value.String() == "true" {
				fmt.Println("Goodbye, have a great day!")
			} else {
				fmt.Println("Bye!")
			}
		},
		Help: "The `goodbye` command prints a farewell message. You can use the `--formal` flag for a formal goodbye.",
	}

	cmd.AddCommand(goodbyeCommand)

	// Execute the CLI
	cmd.Execute()
}
