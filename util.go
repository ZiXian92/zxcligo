package zxcligo

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// NewCliApp takes in command configuration and bootstraps it into a new
// command-line interface application instance.
func NewCliApp(commands []Command) Cli {
	cliInstance := cli{
		config: commands,
	}

	return &cliInstance
}

func newContext(cmdStrings []string) (ctx *Context, err error) {
	ctx = new(Context)
	ctx.Options = make(map[string]interface{})
	args := make([]string, 0, 64)
	isNewFlag := true
	var optName string

	// Set command and arguments
	for _, tok := range cmdStrings[1:] {
		// Check if is option
		if name, val, err := parseOption(tok); err == nil {
			// If previous flag has not been closed with a value efore encountering
			// this flag, we treat it as a boolean flag first.
			// The caller will be responsible for checking against the command's options.
			if !isNewFlag {
				ctx.Options[optName] = "true"
			}

			// Either long name or short name will be true
			if len(val) > 0 {
				ctx.Options[name] = val
				isNewFlag = true
			} else {
				optName = name
				isNewFlag = false
			}
		} else if !isNewFlag {
			// Encountered value that can close the flag
			ctx.Options[optName] = tok
			isNewFlag = true
		} else {
			// Is an argument
			args = append(args, tok)
		}
	}

	// Set command and arguments
	if numArgs := len(args); numArgs > 0 {
		ctx.Cmd = args[0]
		if numArgs > 1 {
			ctx.Args = args[1:]
		}
	} else {
		ctx = nil
		err = fmt.Errorf("no command given")
	}

	// Return completed context and/or error
	return
}

func isLongNameOption(s string) bool {
	res, err := regexp.MatchString("^--[\\w]+(=.+)?$", s)
	if err != nil {
		log.Println("Error matching long name option")
		return false
	}
	return res
}

func isShortNameOption(s string) bool {
	res, err := regexp.MatchString("^-[\\w]+(=.+)?$", s)
	if err != nil {
		log.Println("Error matching short name option")
		return false
	}
	return res
}

// Parses
func parseOption(optString string) (optName, optValue string, err error) {
	// Prepare the option prefix to trim away
	var prefix string
	if isLongNameOption(optString) {
		prefix = "--"
	} else if isShortNameOption(optString) {
		prefix = "-"
	} else {
		err = fmt.Errorf("Invalid option: %s", optString)
		return
	}

	// Trim prefix then try to split into opt-value pair
	tokens := strings.SplitN(strings.TrimPrefix(optString, prefix), "=", 2)
	optName = tokens[0]
	if len(tokens) > 1 {
		optValue = tokens[1]
	} else {
		optValue = ""
	}
	err = nil

	return
}

func findCommand(commands []Command, cmd string) *Command {
	for _, command := range commands {
		if command.Name == cmd {
			return &command
		}
	}
	return nil
}

// Uses the options list to do further processing on the context's options mapping
// so that it can be easily referred to by the user using the option's long name.
func processCommandOptions(ctx *Context, options []Option) error {
	// opts store the processed options.
	// This will replace the current one on ctx.
	opts := make(map[string]interface{})

	for _, opt := range options {
		// Decide on which name the user should use to refer to this
		// option's value.
		var name string
		if len(opt.LongName) > 0 {
			name = opt.LongName
		} else {
			name = opt.ShortName
		}

		// Not sure whether user uses the long name or the short name.
		// So we try both.
		longOptValue, hasLongOptValue := ctx.Options[opt.LongName]
		shortOptValue, hasShortOptValue := ctx.Options[opt.ShortName]

		// The 4 cases of trying to get opt value using long and short names.
		if !hasLongOptValue && !hasShortOptValue {
			// No value supplied. Need to fall back to default value.
			if opt.DefaultValue == nil {
				return fmt.Errorf("option %s is required", name)
			}
			opts[name] = opt.DefaultValue
			continue
		} else if hasLongOptValue && hasShortOptValue {
			// Reject if user uses both name forms of the option
			// due to ambiguity.
			return fmt.Errorf("Conflicting options --%s and -%s given", opt.LongName, opt.ShortName)
		}

		// This variable is to generalize logic in the next part without caring about
		// whether long or short name value is used.
		var optVal string
		if hasLongOptValue {
			optVal = longOptValue.(string)
		} else {
			optVal = shortOptValue.(string)
		}

		// Option value type parsing
		var realOptVal interface{}
		var err error
		switch opt.Type {
		case OptTypeBool:
			realOptVal, err = strconv.ParseBool(optVal)
		case OptTypeFloat:
			realOptVal, err = strconv.ParseFloat(optVal, 64)
		case OptTypeInt:
			realOptVal, err = strconv.ParseInt(optVal, 10, 64)
		case OptTypeString:
			realOptVal, err = optVal, nil
		case OptTypeUint:
			realOptVal, err = strconv.ParseUint(optVal, 10, 64)
		default:
			return fmt.Errorf("Invalid option type for %s", opt.LongName)
		}
		if err != nil {
			return fmt.Errorf("Invalid value type for option %s", opt.LongName)
		}

		// Passed the parsing test, set the option value.
		opts[name] = realOptVal

		// Clean up current options mapping so we can check if there
		// are unexpected options given
		delete(ctx.Options, opt.LongName)
		delete(ctx.Options, opt.ShortName)
	}

	// Test for unexpected options
	for optName := range ctx.Options {
		return fmt.Errorf("Unexpected option %s", optName)
	}

	// Replace current context's options with the processed one.
	ctx.Options = opts

	return nil
}
