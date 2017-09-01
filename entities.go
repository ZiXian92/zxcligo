package zxcligo

// OptType represents the option's value type
type OptType int

const (
	// OptTypeBool is the value to use to indicate boolean type option flag
	OptTypeBool OptType = iota

	// OptTypeFloat is the value to use to indicate float type option flag
	OptTypeFloat

	// OptTypeInt is the value to use to indicate integer type option flag
	OptTypeInt

	// OptTypeString is the value to use to indicate string type option flag
	OptTypeString

	// OptTypeUint is the value to use to indicate unsigned integer type option flag
	OptTypeUint
)

// Option defines a command-line option
type Option struct {
	// Sets the option's value if not provided by user.
	// If not set, it becomes compulsory for the user to pass this option.
	DefaultValue interface{}

	// At least either the LongName or ShortName fields must not be empty.
	// Otherwise, this option will be unusable.
	LongName  string
	ShortName string

	// The value type of this option
	Type OptType

	// Description about this option
	Usage string
}

// Command represents a CLI command for the program,
// together with the option flags supplied to the command.
type Command struct {
	// Name of the command.
	Name string

	// Description of what this command does.
	Usage string

	// The options required by this command.
	Options []Option

	// The function that executes this command.
	Action func(ctx *Context)
}

// Context defines the entitiy that contains parsed command and options
type Context struct {
	// The command passed to the program.
	Cmd string

	// The arguments for the command.
	Args []string

	// The options supplied to the command.
	Options map[string]interface{}
}

// Cli defines the general behaviour of a command-line application
type Cli interface {
	// Runs the command-line program
	Run(cmdStrings []string)
}

type cli struct {
	config []Command
}
