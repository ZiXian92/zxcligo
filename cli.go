package zxcligo

import (
	"fmt"
	"log"
	"os"
)

// NewCliApp takes in command configuration and bootstraps it into a new
// command-line interface application instance.
func NewCliApp() Cli {
	return Cli{}
}

// Cli defines the CLI application instance,
// representing the program that uses this library to bootstrap.
type Cli struct {
	Name     string
	Usage    string
	Version  string
	Flags    []Flag
	Commands []Command
	Action   func(*Context) error
}

// Run takes in command-line arguments list and starts the CLI application.
// Typically os.Args is passed as parameter
func (app *Cli) Run(cmdStrings []string) {
	flagSet := initFlagSet(app.Name, app.Flags)

	// Silence the default behaviour of printing help text
	// when flag parsing fails
	// flagSet.SetOutput(ioutil.Discard)

	// Create new context that we can use to print own formatted help message
	ctx := newContext(app.Name, flagSet, nil)

	// Use of flag.Usage is already called when Parse fails.
	err := flagSet.Parse(cmdStrings[1:])
	if err != nil {
		// printHelp(ctx)
		os.Exit(1)
	}

	ctx.Args = flagSet.Args()

	// If program has no commands configured,
	// simply run its configured action.
	if !app.hasCommands() {
		if app.Action != nil {
			// Running the program results in error.
			if err := app.Action(ctx); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Program logic not implemented.")
		}

		os.Exit(0)
	}

	if len(ctx.Args) == 0 {
		log.Println("Has commands registered but no command given")
		printHelp(ctx)

		os.Exit(1)
	}

	command := app.command(ctx.Args[0])

	// If no match, error. Show help text.
	if command == nil {
		// TODO: Show help text

		os.Exit(1)
	}

	// Else, run the command, passing in ctx.
	command.Run(ctx)
}

func (app *Cli) hasCommands() bool {
	return len(app.Commands) > 0
}

// Gets the command with the given name
func (app *Cli) command(cmd string) *Command {
	for i, command := range app.Commands {
		if command.Name == cmd {
			return &app.Commands[i]
		}
	}
	return nil
}
