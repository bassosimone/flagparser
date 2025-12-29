//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/option.go
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/optiontype.go
//

package flagparser

// Option specifies the kind of option to parse.
type Option struct {
	// DefaultValue is the default value assigned to the option
	// [Value] when the option argument is optional.
	DefaultValue string

	// Prefix is the prefix to use for parsing this option (e.g., `-`)
	Prefix string

	// Name is the option name without the prefix (e.g., `f`).
	Name string

	// Type is the option type.
	Type OptionType
}

// NewOptionWithArgumentNone creates options with no arguments using GNU
// prefixes (- for short, -- for long).
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method
// returns a nil slice.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
func NewOptionWithArgumentNone(shortName byte, longName string) []*Option {
	return newOptionSlice(
		newShortOption(shortName, OptionTypeGroupableArgumentNone),
		newLongOption(longName, OptionTypeStandaloneArgumentNone),
	)
}

// NewEarlyOption creates early options with no arguments using GNU prefixes
// (- for short, -- for long).
//
// You typically use "early" options to register `-h` and `--help` such that
// when the user uses those flags, regardless of whether the command line is
// correct, they see the help text in the output rather than parsing errors.
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method
// returns a nil slice.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
func NewEarlyOption(shortName byte, longName string) []*Option {
	return newOptionSlice(
		newShortOption(shortName, OptionTypeEarlyArgumentNone),
		newLongOption(longName, OptionTypeEarlyArgumentNone),
	)
}

// NewOptionWithArgumentRequired creates options with a required argument using
// GNU prefixes (- for short, -- for long).
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method
// returns a nil slice.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
func NewOptionWithArgumentRequired(shortName byte, longName string) []*Option {
	return newOptionSlice(
		newShortOption(shortName, OptionTypeGroupableArgumentRequired),
		newLongOption(longName, OptionTypeStandaloneArgumentRequired),
	)
}

// NewLongOptionWithArgumentOptional creates a long option with an optional
// argument and a default value using the GNU `--` prefix.
//
// If the long option name is empty, this method returns a nil slice.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
func NewLongOptionWithArgumentOptional(longName, defaultValue string) []*Option {
	if longName == "" {
		return nil
	}
	return []*Option{
		{
			DefaultValue: defaultValue,
			Prefix:       "--",
			Name:         longName,
			Type:         OptionTypeStandaloneArgumentOptional,
		},
	}
}

func newShortOption(shortName byte, optionType OptionType) *Option {
	if shortName == 0 {
		return nil
	}
	return &Option{
		Prefix: "-",
		Name:   string(shortName),
		Type:   optionType,
	}
}

func newLongOption(longName string, optionType OptionType) *Option {
	if longName == "" {
		return nil
	}
	return &Option{
		Prefix: "--",
		Name:   longName,
		Type:   optionType,
	}
}

func newOptionSlice(options ...*Option) []*Option {
	var out []*Option
	for _, option := range options {
		if option != nil {
			out = append(out, option)
		}
	}
	return out
}

// OptionType is the type of an [Option].
type OptionType int64

const (
	optionKindEarly = OptionType(1 << (iota + 4))
	optionKindStandalone
	optionKindGroupable
)

const (
	optionArgumentNone = OptionType(1 << iota)
	optionArgumentRequired
	optionArgumentOptional
)

func (ot OptionType) isEarly() bool {
	return (ot & optionKindEarly) != 0
}

func (ot OptionType) isStandalone() bool {
	return (ot & optionKindStandalone) != 0
}

func (ot OptionType) isGroupable() bool {
	return (ot & optionKindGroupable) != 0
}

// These constants define the allowed [OptionType] values.
const (
	// OptionTypeEarlyArgumentNone indicates an early option requiring no arguments.
	//
	// Typically used for `-h` and `--help`.
	OptionTypeEarlyArgumentNone = optionKindEarly | optionArgumentNone

	// OptionTypeStandaloneArgumentNone indicates a standalone option requiring no arguments.
	//
	// Typically used for `--verbose` or `--quiet`.
	OptionTypeStandaloneArgumentNone = optionKindStandalone | optionArgumentNone

	// OptionTypeStandaloneArgumentRequired indicates a standalone option requiring an argument.
	//
	// Typically used for stuff like `--output FILE` (or `--output=FILE`).
	OptionTypeStandaloneArgumentRequired = optionKindStandalone | optionArgumentRequired

	// OptionTypeStandaloneArgumentOptional indicates a standalone option with an optional argument.
	//
	// Typically used for stuff like `--http=1.1` (or `--http` to get the default).
	OptionTypeStandaloneArgumentOptional = optionKindStandalone | optionArgumentOptional

	// OptionTypeGroupableArgumentNone indicates a groupable option requiring no arguments.
	//
	// Typically used for options like `-v` (for verbose).
	//
	// These options can be grouped together like in `-xvzd DIR`.
	OptionTypeGroupableArgumentNone = optionKindGroupable | optionArgumentNone

	// OptionTypeGroupableArgumentRequired indicates groupable option requiring an argument.
	//
	// Typically used for options like `-d DIR` or `-dDIR` (to select a directory).
	//
	// These options can be grouped together like in `-xvzd DIR`.
	OptionTypeGroupableArgumentRequired = optionKindGroupable | optionArgumentRequired
)
