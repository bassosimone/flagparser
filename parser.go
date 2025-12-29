//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/parser.go
//

package flagparser

import (
	"fmt"

	"github.com/bassosimone/flagscanner"
	"github.com/bassosimone/runtimex"
)

// Parser is a command line parser.
//
// Construct with [NewParser] to get GNU parsing semantics. Otherwise, if you
// need a distinct parser semantics, please construct manually.
type Parser struct {
	// DisablePermute optionally disables permuting options and arguments.
	//
	// Consider the following command line arguments:
	//
	// 	https://www.google.com/ -H 'Host: google.com'
	//
	// The default behavior is to permute this to:
	//
	// 	-H 'Host: google.com' https://www.google.com/
	//
	// However, when DisablePermute is true, we keep the command
	// line unmodified. While permuting is a nice-to-have property
	// in general, consider instead the following case:
	//
	// 	foreach -kx git status -v
	//
	// With permutation, this command line would become:
	//
	// 	-kv -v foreach git status
	//
	// This is not the desired behavior if the foreach subcommand
	// takes another command and its options as arguments.
	//
	// To make the above command line work with permutation, a
	// user would instead need to write this:
	//
	// 	foreach -kx -- git status -v
	//
	// By setting DisablePermute to true, the `--` separator
	// becomes unnecessary and the UX is improved.
	DisablePermute bool

	// MaxPositionalArguments is the maximum number of positional
	// arguments allowed by the parser. The default is zero, meaning
	// that the parser won't accept more than zero positionals.
	MaxPositionalArguments int

	// MinPositionalArguments is the minimum number of positional
	// arguments allowed by the parser. The default is zero, meaning
	// that the parser won't accept less than zero positionals.
	MinPositionalArguments int

	// OptionsArgumentsSeparator is the optional separator that terminates
	// the parsing of options, treating all remaining tokens in the command
	// line as positional arguments. The default is empty, meaning that
	// the parser will always parse all the available options.
	OptionsArgumentsSeparator string

	// Options contains the optional options configured for this parser.
	//
	// When parsing, we will ensure there are no duplicate option names or
	// ambiguous separators across all options.
	//
	// If you don't set this field, the parser will automatically
	// configure itself to parse GNU-style options, meaning that it
	// will use `-` as the prefix for short options and `--` as
	// the prefix for long options. No options will be defined so
	// any option will be considered unknown and cause a parse error.
	Options []*Option
}

// ErrTooFewPositionalArguments is returned when the number of positional
// arguments is less than the configured minimum.
type ErrTooFewPositionalArguments struct {
	// Min is the minimum number of positional arguments required.
	Min int

	// Have is the number of positional arguments provided.
	Have int
}

var _ error = ErrTooFewPositionalArguments{}

// Error returns a string representation of this error.
func (err ErrTooFewPositionalArguments) Error() string {
	return fmt.Sprintf("too few positional arguments: expected at least %d, got %d", err.Min, err.Have)
}

// ErrTooManyPositionalArguments is returned when the number of positional
// arguments is greater than the configured maximum.
type ErrTooManyPositionalArguments struct {
	// Max is the maximum number of positional arguments allowed.
	Max int

	// Have is the number of positional arguments provided.
	Have int
}

var _ error = ErrTooManyPositionalArguments{}

// Error returns a string representation of this error.
func (err ErrTooManyPositionalArguments) Error() string {
	return fmt.Sprintf("too many positional arguments: expected at most %d, got %d", err.Max, err.Have)
}

// NewParser creates a new [*Parser] following the GNU convention.
//
// Specifically, we use these settings:
//
//  1. command line permutation is enabled
//
//  2. zero positional arguments are allowed
//
//  3. the separator is set to `--`
//
//  4. no options have been defined yet
//
// Create [*Parser] manually when you need different defaults.
func NewParser() *Parser {
	return &Parser{
		DisablePermute:            false,
		MaxPositionalArguments:    0,
		MinPositionalArguments:    0,
		OptionsArgumentsSeparator: "--",
		Options:                   []*Option{},
	}
}

// SetMinMaxPositionalArguments sets the minimum and maximum positional arguments.
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid minimum and maximum positional arguments values will cause
// no errors until you attempt to parse the command line.
func (px *Parser) SetMinMaxPositionalArguments(minArgs, maxArgs int) {
	px.MinPositionalArguments = minArgs
	px.MaxPositionalArguments = maxArgs
}

// AddOption adds one or more options to the parser.
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
func (px *Parser) AddOption(options ...*Option) {
	for _, option := range options {
		if option != nil {
			px.Options = append(px.Options, option)
		}
	}
}

// AddOptionWithArgumentNone adds a short and long option taking no argument
// and using the `-` and `--` prefixes, which follow the GNU conventions.
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method is
// a no-operation that does not change the [*Parser].
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
//
// Use [NewOptionWithArgumentNone] to construct options without mutating the parser.
func (px *Parser) AddOptionWithArgumentNone(shortName byte, longName string) {
	px.AddOption(NewOptionWithArgumentNone(shortName, longName)...)
}

// AddEarlyOption adds an "early" short and long option taking no argument and
// using the `-` and `--` prefixes, which follow the GNU conventions.
//
// You typically use "early" options to register `-h` and `--help` such that when
// the user uses those flags, regardless of whether the command line is correct, they
// see the help text in the output rather than parsing errors.
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method is
// a no-operation that does not change the [*Parser].
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
//
// Use [NewEarlyOption] to construct options without mutating the parser.
func (px *Parser) AddEarlyOption(shortName byte, longName string) {
	px.AddOption(NewEarlyOption(shortName, longName)...)
}

// AddOptionWithArgumentRequired adds a short and long option with a required argument
// and using the `-` and `--` prefixes, which follow the GNU conventions.
//
// A zero short option value skips adding the short option. An empty long option
// value skips adding the long option. If both are zero/empty, this method is
// a no-operation that does not change the [*Parser].
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
//
// Use [NewOptionWithArgumentRequired] to construct options without mutating the parser.
func (px *Parser) AddOptionWithArgumentRequired(shortName byte, longName string) {
	px.AddOption(NewOptionWithArgumentRequired(shortName, longName)...)
}

// AddLongOptionWithArgumentOptional adds a long option with an optional argument
// with the given default value and using `--` prefix, which follows the GNU conventions.
//
// If the long option name is empty, this method is a no-operation that does
// not change the [*Parser].
//
// This method MUTATES [*Parser] and is NOT SAFE to call concurrently.
//
// Setting invalid option names (e.g., a duplicate option name) will cause
// no errors until you attempt to parse the command line.
//
// Use [NewLongOptionWithArgumentOptional] to construct options without mutating the parser.
func (px *Parser) AddLongOptionWithArgumentOptional(longName, defaultValue string) {
	px.AddOption(NewLongOptionWithArgumentOptional(longName, defaultValue)...)
}

// Parse parses the command line arguments.
//
// This method does not mutate [*Parser] and is safe to call concurrently.
//
// The args MUST NOT include the program name.
func (px *Parser) Parse(args []string) ([]Value, error) {
	// Create the configuration
	cfg, err := newConfig(px)
	if err != nil {
		return nil, err
	}

	// Create scanner for the parser.
	sx := &flagscanner.Scanner{
		Separator: px.OptionsArgumentsSeparator,
		Prefixes:  []string{},
	}
	for prefix := range cfg.prefixes {
		sx.Prefixes = append(sx.Prefixes, prefix)
	}

	// Tokenize the command line arguments.
	tokens := sx.Scan(args)

	// Preflight the command line arguments searching for early options
	// and return immediately when found. This algorithm allows for
	// immediately intercepting `--help` regardless of possibly invalid
	// options, which, in turn, improves the UX, because we can show
	// the full help to the user rather than errors.
	if value, found := earlyParse(px.Options, tokens); found {
		return []Value{value}, nil
	}

	// Create a deque with the values to parse.
	input := &deque[flagscanner.Token]{values: tokens}

	// Parse the command line.
	var (
		options     = &deque[Value]{}
		positionals = &deque[Value]{}
	)
	if err := doParse(cfg, input, options, positionals); err != nil {
		return nil, err
	}

	// Ensure this stage has emptied the input deque.
	runtimex.Assert(input.Empty())

	// Ensure the number of positional arguments is within the limits.
	if len(positionals.values) < px.MinPositionalArguments {
		return nil, ErrTooFewPositionalArguments{
			Min:  px.MinPositionalArguments,
			Have: len(positionals.values),
		}
	}
	if len(positionals.values) > px.MaxPositionalArguments {
		return nil, ErrTooManyPositionalArguments{
			Max:  px.MaxPositionalArguments,
			Have: len(positionals.values),
		}
	}

	// Create the result slice by optionally permuting the entries.
	result := permute(cfg.disablePermute(), options.values, positionals.values)
	return result, nil
}
