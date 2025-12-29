//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/parse.go
//

package flagparser

import (
	"fmt"
	"io"
	"strings"

	"github.com/bassosimone/flagscanner"
)

// ErrOptionRequiresNoArgument indicates that an argument was
// passed to an option that requires no arguments.
type ErrOptionRequiresNoArgument struct {
	// Option is the offending option
	Option *Option

	// Token is the related token
	Token flagscanner.Token
}

var _ error = ErrOptionRequiresNoArgument{}

// Error returns a string representation of this error.
func (err ErrOptionRequiresNoArgument) Error() string {
	return fmt.Sprintf("option requires no argument: %s%s", err.Option.Prefix, err.Option.Name)
}

// ErrOptionRequiresArgument indicates that no argument was
// passed to an option that requires an argument.
type ErrOptionRequiresArgument struct {
	// Option is the offending option
	Option *Option

	// Token is the related token
	Token flagscanner.Token
}

var _ error = ErrOptionRequiresArgument{}

// Error returns a string representation of this error.
func (err ErrOptionRequiresArgument) Error() string {
	return fmt.Sprintf("option requires an argument: %s%s", err.Option.Prefix, err.Option.Name)
}

// parseDebugWriter is only used by tests to surface parsing steps.
var parseDebugWriter = io.Discard

func doParse(cfg *config, input *deque[flagscanner.Token], options, positionals *deque[Value]) error {
	// Know when to treat everything else as positional
	var onlypositionals bool

	// Attempt to consume all the available tokens
	for !input.Empty() {
		// Get the current token and advance
		cur, _ := input.Front()
		input.PopFront()
		fmt.Fprintf(parseDebugWriter, "\nprocessing token: %+v\n", cur)

		// Decide what to do depending on the token type
		switch cur := cur.(type) {

		// On positional argument, stop parsing if permutation is disabled
		case flagscanner.PositionalArgumentToken:
			value := ValuePositionalArgument{
				Tok:   cur,
				Value: cur.Value,
			}
			positionals.PushBack(value)
			fmt.Fprintf(parseDebugWriter, "added positional argument value: %+v\n", value)
			if cfg.disablePermute() {
				fmt.Fprint(parseDebugWriter, "no permute: starting to treat everything as positional\n")
				onlypositionals = true
			}
			continue

		// Stop parsing if we encounter the options-arguments separator
		case flagscanner.OptionsArgumentsSeparatorToken:
			value := ValueOptionsArgumentsSeparator{
				Tok:       cur,
				Separator: cur.Separator,
			}
			positionals.PushBack(value)
			fmt.Fprintf(parseDebugWriter, "added options-arguments separator value: %+v\n", value)
			fmt.Fprint(parseDebugWriter, "seen separator: starting to treat everything as positional\n")
			onlypositionals = true
			continue

		// OK, we've got an option, we're definitely interested
		case flagscanner.OptionToken:
			// When we're treating everything as positional, just short-circuit it
			if onlypositionals {
				value := ValuePositionalArgument{
					Tok:   cur,
					Value: cur.String(),
				}
				positionals.PushBack(value)
				fmt.Fprintf(parseDebugWriter, "added option as positional value: %+v\n", value)
				continue
			}

			// Switch on the kind of flag based on standalone vs groupable
			optkind := cfg.prefixes[cur.Prefix]
			switch {
			case optkind.isStandalone():
				if err := doParseStandaloneOption(cfg, cur, input, options); err != nil {
					return err
				}

			case optkind.isGroupable():
				if err := doParseGroupableOption(cfg, cur, input, options); err != nil {
					return err
				}

			default:
				panic(fmt.Sprintf("unhandled option type: %d", optkind))
			}
		}
	}
	return nil
}

func doParseStandaloneOption(
	cfg *config, cur flagscanner.OptionToken, input *deque[flagscanner.Token], options *deque[Value]) error {
	// The option may contain a value, account for this
	var optname, optvalue string
	index := strings.Index(cur.Name, "=")
	if index > 0 {
		optname = cur.Name[:index]
		optvalue = cur.Name[index+1:]
	} else {
		optname = cur.Name
	}
	fmt.Fprintf(parseDebugWriter, "optname=%q, optvalue=%q\n", optname, optvalue)

	// Obtain the option given its name and prefix
	option, err := cfg.findOption(cur, optname, optionKindStandalone)
	if err != nil {
		return err
	}
	fmt.Fprintf(parseDebugWriter, "found option: %+v\n", option)

	// Specialize handling depending on the option type
	switch option.Type {
	case OptionTypeStandaloneArgumentNone:
		if optname != cur.Name { // account for `--option=VALUE` case
			return ErrOptionRequiresNoArgument{Option: option, Token: cur}
		}

	case OptionTypeStandaloneArgumentOptional:
		if optvalue == "" {
			optvalue = option.DefaultValue
		}

	case OptionTypeStandaloneArgumentRequired:
		if optname == cur.Name { // account for `--option VALUE` case
			if input.Empty() {
				return ErrOptionRequiresArgument{Option: option, Token: cur}
			}
			tok, _ := input.Front()
			input.PopFront()
			optvalue = tok.String()
		}

	default:
		panic(fmt.Sprintf("unhandled option type: %d", option.Type))
	}

	// Create and add the option
	value := ValueOption{Option: option, Tok: cur, Value: optvalue}
	options.PushBack(value)
	fmt.Fprintf(parseDebugWriter, "added option value: %+v\n", value)
	return nil
}

func doParseGroupableOption(
	cfg *config, cur flagscanner.OptionToken, input *deque[flagscanner.Token], options *deque[Value]) error {
	// Scan through each byte inside the option group
	for otokname := cur.Name; len(otokname) > 0; {
		// Extract the option name and advance
		optname := otokname[0]
		otokname = otokname[1:]
		fmt.Fprintf(parseDebugWriter, "optname=%q\n", string(optname))

		// Obtain the option given its name and prefix
		option, err := cfg.findOption(cur, string(optname), optionKindGroupable)
		if err != nil {
			return err
		}
		fmt.Fprintf(parseDebugWriter, "found option: %+v\n", option)

		// Specialize handling depending on option type
		var optvalue string
		switch option.Type {
		case OptionTypeGroupableArgumentNone:
			// nothing

		case OptionTypeGroupableArgumentRequired:
			switch {
			case len(otokname) > 0: // the `-vfFILE` GNU-compatible case
				optvalue = otokname
				otokname = ""

			case !input.Empty(): // the `-vf FILE` case
				tok, _ := input.Front()
				input.PopFront()
				optvalue = tok.String()

			default:
				return ErrOptionRequiresArgument{Option: option, Token: cur}
			}

		default:
			panic(fmt.Sprintf("unhandled option type: %d", option.Type))
		}

		// Create and add the option
		value := ValueOption{Option: option, Tok: cur, Value: optvalue}
		options.PushBack(value)
		fmt.Fprintf(parseDebugWriter, "added option value: %+v\n", value)
	}
	return nil
}
