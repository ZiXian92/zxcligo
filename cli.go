package zxcligo

import (
	"fmt"
	"os"
)

// Run takes in command-line arguments list and starts the CLI application.
// Typically os.Args is passed as parameter
func (app *cli) Run(cmdStrings []string) {
	ctx, err := newContext(cmdStrings)
	if err != nil {
		fmt.Printf("%v\n", err)

		// TODO: Show help text

		os.Exit(1)
	}

	command := findCommand(app.config, ctx.Cmd)
	if command == nil {
		fmt.Printf("Invalid command %s\n", ctx.Cmd)

		// TODO: Show help text

		os.Exit(1)
	}

	err = processCommandOptions(ctx, command.Options)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if command.Action == nil {
		fmt.Printf("Command %s not implemented\n", command.Name)
		os.Exit(1)
	}

	// Execute the action associated with the command
	command.Action(ctx)
}
