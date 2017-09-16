package zxcligo

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Returns the actual value that is passed as a Flag.
// Pointer values passed s Flag are dereferenced when returned.
func getFlagValue(f Flag) reflect.Value {
	v := reflect.ValueOf(f)

	if v.Kind() == reflect.Ptr {
		return reflect.Indirect(v)
	}

	return v
}

// Interface definition to sort array of strings by string length.
type byLength []string

func (s byLength) Len() int {
	return len(s)
}

func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

// For returns [short_flag] [placeholder], [long_flag] [placeholder]
func formattedFlagUsagePrefix(flagName, placeholder string) string {
	tokens := strings.Split(flagName, ",")
	for i, token := range tokens {
		tokens[i] = strings.TrimSpace(token)
	}

	sort.Sort(byLength(tokens))

	// After sorting tokens by length, we make it such that the very first token
	// is the flag to be used with - prefix, others with -- prefix.
	var prefixBuffer bytes.Buffer
	numTokens := len(tokens)
	placeholderIsEmpty := len(placeholder) == 0
	for i, token := range tokens {
		if i == 0 {
			prefixBuffer.WriteString("-")
		} else {
			prefixBuffer.WriteString("--")
		}

		prefixBuffer.WriteString(token)

		if !placeholderIsEmpty {
			prefixBuffer.WriteString(" <")
			prefixBuffer.WriteString(placeholder)
			prefixBuffer.WriteString(" >")
		}

		if i < numTokens-1 {
			prefixBuffer.WriteString(", ")
		}
	}

	return prefixBuffer.String()
}

// stringifyFlag returns a formatted string form of the given flag
// for usage display purpose.
func stringifyFlag(f Flag) string {
	actualFlag := getFlagValue(f)

	valPlaceholder := ""
	if ph := actualFlag.FieldByName("Placeholder"); ph.IsValid() {
		valPlaceholder = ph.String()
	}

	prefix := formattedFlagUsagePrefix(actualFlag.FieldByName("Name").String(), valPlaceholder)
	usage := actualFlag.FieldByName("Usage").String()

	return fmt.Sprintf("%s\t%s", prefix, usage)
}

// Splits flagNames into names using known name delimiter ", "
// and then applies the registerFn for each name.
func registerFlagForAllNames(flagNames string, registerFn func(string)) {
	names := strings.Split(flagNames, ",")
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}

	for _, name := range names {
		registerFn(name)
	}
}

// Applies the
func initFlagSet(appName string, flags []Flag) *flag.FlagSet {
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	// Apply each flag to the flagset
	for _, flag := range flags {
		flag.Apply(flagSet)
	}

	return flagSet
}
