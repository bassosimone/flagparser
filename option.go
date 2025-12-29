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
