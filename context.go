package zxcligo

import (
	"flag"
	"strconv"
)

// Context defines the entitiy that contains parsed command and options
type Context struct {
	Command string
	Args    []string
	flags   *flag.FlagSet
	Parent  *Context
}

func newContext(name string, flagSet *flag.FlagSet, parentCtx *Context) *Context {
	return &Context{
		Command: name,
		Args:    flagSet.Args(),
		flags:   flagSet,
		Parent:  parentCtx,
	}
}

// BoolFlag returns the flag value corresponding to the flag
// with the given name.
func (ctx *Context) BoolFlag(name string) (val bool, ok bool) {
	flg := ctx.flags.Lookup(name)

	if flg == nil {
		// Flag not found, check if it is some global flag
		if ctx.Parent != nil {
			return ctx.Parent.BoolFlag(name)
		}
		ok = false
		return
	}

	val, err := strconv.ParseBool(flg.Value.String())
	if err != nil {
		ok = false
	} else {
		ok = true
	}

	return
}

// Float64Flag returns the flag value of the float flag
// with the given name.
func (ctx *Context) Float64Flag(name string) (val float64, ok bool) {
	f := ctx.flags.Lookup(name)

	if f == nil {
		// Flag not found, check if it is some global flag
		if ctx.Parent != nil {
			return ctx.Parent.Float64Flag(name)
		}
		ok = false
		return
	}

	val, err := strconv.ParseFloat(f.Value.String(), 64)
	if err != nil {
		ok = false
	} else {
		ok = true
	}

	return
}

// Int64Flag returns the flag value of the integer flag
// with the given name.
func (ctx *Context) Int64Flag(name string) (val int64, ok bool) {
	f := ctx.flags.Lookup(name)

	if f == nil {
		// Flag not found, check if it is some global flag
		if ctx.Parent != nil {
			return ctx.Parent.Int64Flag(name)
		}
		ok = false
		return
	}

	val, err := strconv.ParseInt(f.Value.String(), 10, 64)
	if err != nil {
		ok = false
	} else {
		ok = true
	}

	return
}

// StringFlag returns the flag value of the string flag
// with the given name.
func (ctx *Context) StringFlag(name string) (val string, ok bool) {
	f := ctx.flags.Lookup(name)

	if f == nil {
		// Flag not found, check if it is some global flag
		if ctx.Parent != nil {
			return ctx.Parent.StringFlag(name)
		}
		ok = false
		return
	}

	return f.Value.String(), true
}

// Uint64Flag returns the flag value of the unsigned integer flag
// with the given name.
func (ctx *Context) Uint64Flag(name string) (val uint64, ok bool) {
	f := ctx.flags.Lookup(name)

	if f == nil {
		// Flag not found, check if it is some global flag
		if ctx.Parent != nil {
			return ctx.Parent.Uint64Flag(name)
		}
		ok = false
		return
	}

	val, err := strconv.ParseUint(f.Value.String(), 10, 64)
	if err != nil {
		ok = false
	} else {
		ok = true
	}

	return
}
