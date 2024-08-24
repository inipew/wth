package main

import (
	"cmd/internal"
	"flag"
	"fmt"
)

func main() {
	cli := internal.NewCLI("1.0.0")

	// Membuat command `hello`
	helloFlags := &internal.FlagSetParser{flag.NewFlagSet("hello", flag.ContinueOnError)}
	helloFlags.String("name", "World", "Name to greet")

	helloCommand := &internal.Command{
		Name:        "hello",
		Description: "Prints a greeting message.",
		Flags:       helloFlags,
		Run: func(args []string) {
			name := helloFlags.Lookup("name").Value.String()
			fmt.Printf("Hello, %s!\n", name)
		},
		Help: "The `hello` command prints a greeting message. You can use the `--name` flag to specify the name.",
	}

	cli.AddCommand(helloCommand)

	// Membuat command `goodbye`
	goodbyeFlags := &internal.FlagSetParser{flag.NewFlagSet("goodbye", flag.ContinueOnError)}
	goodbyeFlags.Bool("formal", false, "Use formal goodbye")

	goodbyeCommand := &internal.Command{
		Name:        "goodbye",
		Description: "Prints a farewell message.",
		Flags:       goodbyeFlags,
		Run: func(args []string) {
			if goodbyeFlags.Lookup("formal").Value.String() == "true" {
				fmt.Println("Goodbye, have a great day!")
			} else {
				fmt.Println("Bye!")
			}
		},
		Help: "The `goodbye` command prints a farewell message. You can use the `--formal` flag for a formal goodbye.",
	}

	cli.AddCommand(goodbyeCommand)

	// Menjalankan CLI
	cli.Execute()
}
