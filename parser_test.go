//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/parser_test.go
//

package flagparser

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrTooFewPositionalArguments(t *testing.T) {
	err := ErrTooFewPositionalArguments{Min: 3, Have: 1}
	want := "too few positional arguments: expected at least 3, got 1"
	assert.Equal(t, want, err.Error())
}

func TestErrTooManyPositionalArguments(t *testing.T) {
	err := ErrTooManyPositionalArguments{Max: 3, Have: 5}
	want := "too many positional arguments: expected at most 3, got 5"
	assert.Equal(t, want, err.Error())
}

func TestParser_AddOption(t *testing.T) {
	px := NewParser()
	option := &Option{Prefix: "--", Name: "verbose"}
	px.AddOption(nil, option, nil)
	assert.Equal(t, []*Option{option}, px.Options)
}

func TestParser_Parse(t *testing.T) {
	// Note: example_test.go covers many parsing cases; this file focuses on
	// configuration and error paths not easily expressed as examples.
	// Define the test case structure
	type testcase struct {
		args        []string       // argument vector to parse (excluding program name)
		newParser   func() *Parser // return the parser to use
		expectValue []string       // expected parsed and reserialized values
		expectErr   error          // expected error, if any
	}

	cases := []testcase{

		{
			args: []string{"https://example.com/file.txt", "-fsSL", "--remote-name"},
			newParser: func() *Parser {
				px := NewParser()
				px.SetMinMaxPositionalArguments(1, math.MaxInt)
				px.AddOptionWithArgumentNone('f', "fail")
				px.AddOptionWithArgumentNone('L', "location")
				px.AddOptionWithArgumentRequired('O', "remote-name")
				px.AddOptionWithArgumentNone('S', "show-error")
				px.AddOptionWithArgumentNone('s', "silent")
				return px
			},
			expectValue: []string{},
			expectErr:   errors.New("option requires an argument: --remote-name"),
		},

		{
			args: []string{"https://example.com/file.txt", "-fsSL", "--remote-name=FOO"},
			newParser: func() *Parser {
				px := NewParser()
				px.SetMinMaxPositionalArguments(1, math.MaxInt)
				px.AddOptionWithArgumentNone('f', "fail")
				px.AddOptionWithArgumentNone('L', "location")
				px.AddOptionWithArgumentNone('O', "remote-name")
				px.AddOptionWithArgumentNone('S', "show-error")
				px.AddOptionWithArgumentNone('s', "silent")
				return px
			},
			expectValue: []string{},
			expectErr:   errors.New("option requires no argument: --remote-name"),
		},

		{
			args: []string{"@8.8.8.8", "-p53", "IN", "+short", "A", "-h"},
			newParser: func() *Parser {
				return &Parser{
					MinPositionalArguments: 1,
					MaxPositionalArguments: 4,
					Options:                NewEarlyOption('h', ""),
				}
			},
			expectValue: []string{"-h"},
			expectErr:   nil,
		},

		{
			args: []string{"@8.8.8.8", "-P53", "IN", "+short", "A"},
			newParser: func() *Parser {
				return &Parser{
					MinPositionalArguments: 1,
					MaxPositionalArguments: 4,
					Options:                NewOptionWithArgumentRequired('p', ""),
				}
			},
			expectValue: nil,
			expectErr:   errors.New("unknown option: -P"),
		},

		{
			args: []string{"@8.8.8.8", "-p53", "IN", "+short", "A", "-p"},
			newParser: func() *Parser {
				return &Parser{
					MinPositionalArguments: 1,
					MaxPositionalArguments: 4,
					Options: []*Option{
						{
							Name:   "h",
							Prefix: "-",
							Type:   OptionTypeEarlyArgumentNone,
						},
						{
							Name:   "p",
							Prefix: "-",
							Type:   OptionTypeGroupableArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("option requires an argument: -p"),
		},

		{
			args: []string{},
			newParser: func() *Parser {
				return &Parser{
					Options: []*Option{
						{
							Name:   "h",
							Prefix: "-",
							Type:   OptionTypeEarlyArgumentNone,
						},
						{
							Name:   "port",
							Prefix: "-",
							Type:   OptionTypeGroupableArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("groupable option names should be a single byte, found: &{DefaultValue: Prefix:- Name:port Type:66}"),
		},

		{
			args: []string{},
			newParser: func() *Parser {
				return &Parser{
					Options: []*Option{
						{
							Name:   "p",
							Prefix: "--",
							Type:   OptionTypeStandaloneArgumentRequired,
						},
						{
							Name:   "p",
							Prefix: "-",
							Type:   OptionTypeGroupableArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("multiple options with \"p\" name"),
		},

		{
			args: []string{},
			newParser: func() *Parser {
				return &Parser{
					Options: []*Option{
						{
							Name:   "",
							Prefix: "--",
							Type:   OptionTypeStandaloneArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("option name cannot be empty: &{DefaultValue: Prefix:-- Name: Type:34}"),
		},

		{
			args: []string{},
			newParser: func() *Parser {
				return &Parser{
					Options: []*Option{
						{
							Name:   "short",
							Prefix: "",
							Type:   OptionTypeStandaloneArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("option prefix cannot be empty: &{DefaultValue: Prefix: Name:short Type:34}"),
		},

		{
			args: []string{},
			newParser: func() *Parser {
				return &Parser{
					Options: []*Option{
						{
							Name:   "short",
							Prefix: "-",
							Type:   OptionTypeStandaloneArgumentRequired,
						},
						{
							Name:   "p",
							Prefix: "-",
							Type:   OptionTypeGroupableArgumentRequired,
						},
					},
				}
			},
			expectValue: nil,
			expectErr:   errors.New("prefix \"-\" is used for both standalone and groupable options"),
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			// Parse the arguments using the parser
			px := tc.newParser()
			values, err := px.Parse(tc.args)

			// Check for expected error
			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				return
			}
			assert.NoError(t, err)

			// Compare to expectation
			got := []string{}
			for _, entry := range values {
				got = append(got, entry.Strings()...)
			}
			assert.Equal(t, tc.expectValue, got)
		})
	}
}

func TestParserEmptyDefaultsToGNUStyleOptions(t *testing.T) {
	// Create a new empty parser with no options
	px := &Parser{}

	// Parse arguments
	values, err := px.Parse([]string{"--option", "value"})

	// Make sure the error is an unknown option, which means that
	// we were configured to use the GNU style parser
	var unknownOption ErrUnknownOption
	good := errors.As(err, &unknownOption)
	assert.True(t, good)
	assert.Nil(t, values)
}
