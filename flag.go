package zxcligo

import (
	"flag"
	"fmt"
)

// Flag defines the generic behaviour of the various
// types of command-line flags.
type Flag interface {
	fmt.Stringer
	Apply(*flag.FlagSet)
}

// BoolFlag defines a boolean flag such that
// giving this flag in the commandline sets it to true,
// and false otherwise.
type BoolFlag struct {
	Name  string
	Usage string
}

// String returns the string representation of this flag
// for use in printing help message.
func (b BoolFlag) String() string {
	return stringifyFlag(b)
}

// Apply applies this flag to the given flagset
// representing the current commandline context.
func (b BoolFlag) Apply(flagSet *flag.FlagSet) {
	registerFlagForAllNames(b.Name, func(name string) {
		flagSet.Bool(name, false, b.Usage)
	})
}

// Int64Flag defines a signed integer flag
type Int64Flag struct {
	Name        string
	Usage       string
	Value       int64
	Placeholder string
}

// String returns the string representation of this flag
// for use in printing help message.
func (i Int64Flag) String() string {
	return stringifyFlag(i)
}

// Apply applies this flag to the given flagset
// representing the current commandline context.
func (i Int64Flag) Apply(flagSet *flag.FlagSet) {
	registerFlagForAllNames(i.Name, func(name string) {
		flagSet.Int64(name, i.Value, i.Usage)
	})
}

// Uint64Flag defines an unsigned integer flag.
type Uint64Flag struct {
	Name        string
	Usage       string
	Value       uint64
	Placeholder string
}

// String returns the string representation of this flag
// for use in printing help message.
func (u Uint64Flag) String() string {
	return stringifyFlag(u)
}

// Apply applies this flag to the given flagset
// representing the current commandline context.
func (u Uint64Flag) Apply(flagSet *flag.FlagSet) {
	registerFlagForAllNames(u.Name, func(name string) {
		flagSet.Uint64(name, u.Value, u.Usage)
	})
}

// Float64Flag defines a floating-point number flag.
type Float64Flag struct {
	Name        string
	Usage       string
	Value       float64
	Placeholder string
}

// String returns the string representation of this flag
// for use in printing help message.
func (f Float64Flag) String() string {
	return stringifyFlag(f)
}

// Apply applies this flag to the given flagset
// representing the current commandline context.
func (f Float64Flag) Apply(flagSet *flag.FlagSet) {
	registerFlagForAllNames(f.Name, func(name string) {
		flagSet.Float64(name, f.Value, f.Usage)
	})
}

// StringFlag defines a string value flag.
type StringFlag struct {
	Name        string
	Usage       string
	Value       string
	PlaceHolder string
}

// String returns the string representation of this flag
// for use in printing help message.
func (s StringFlag) String() string {
	return stringifyFlag(s)
}

// Apply applies this flag to the given flagset
// representing the current commandline context.
func (s StringFlag) Apply(flagSet *flag.FlagSet) {
	registerFlagForAllNames(s.Name, func(name string) {
		flagSet.String(name, s.Value, s.Usage)
	})
}
