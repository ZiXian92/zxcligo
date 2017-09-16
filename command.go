package zxcligo

import (
	"log"
	"os"
)

// Command represents a CLI command for the program,
// together with the option flags supplied to the command.
type Command struct {
	// Name of the command.
	Name string

	// Description of what this command does.
	Usage string

	// The options required by this command.
	Flags []Flag

	// The function that executes this command.
	Action func(ctx *Context)
}

// Run runs this command with the given command-line context.
func (c *Command) Run(ctx *Context) {
	// Here the condition must be that ctx.Args has at least 1 element.
	flagSet := initFlagSet(c.Name, c.Flags)

	newCtx := newContext(c.Name, flagSet, ctx)
	if err := flagSet.Parse(ctx.Args[1:]); err != nil {
		os.Exit(1)
	}

	newCtx.Args = flagSet.Args()

	if c.Action == nil {
		log.Printf("No registered action")
		os.Exit(1)
	}

	c.Action(newCtx)
}
