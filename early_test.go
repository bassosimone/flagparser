//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/early_test.go
//

package flagparser

import (
	"testing"

	"github.com/bassosimone/flagscanner"
	"github.com/stretchr/testify/assert"
)

// Ensure that the earlyFind algorithm is working as intended.
func Test_earlyFind(t *testing.T) {
	// Define the options we recognize
	options := []*Option{
		{
			Prefix: "-",
			Name:   "h",
			Type:   OptionTypeEarlyArgumentNone,
		},
		{
			Prefix: "--",
			Name:   "help",
			Type:   OptionTypeEarlyArgumentNone,
		},
		{
			Prefix: "+",
			Name:   "short",
			Type:   OptionTypeStandaloneArgumentNone,
		},
	}

	// Define the test cases
	type testcase struct {
		name   string              // name of the test case
		tokens []flagscanner.Token // tokens to parse
		expect Value               // expected parsed value
	}

	cases := []testcase{
		{
			name: "successful recognition of --help",
			tokens: []flagscanner.Token{
				flagscanner.OptionToken{Idx: 0, Prefix: "-", Name: "x"},
				flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
				flagscanner.OptionToken{Idx: 2, Prefix: "--", Name: "help"},
			},
			expect: ValueOption{
				Option: options[1],
				Tok:    flagscanner.OptionToken{Idx: 2, Prefix: "--", Name: "help"},
			},
		},

		{
			name: "successful recognition of -h",
			tokens: []flagscanner.Token{
				flagscanner.OptionToken{Idx: 0, Prefix: "-", Name: "x"},
				flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
				flagscanner.OptionToken{Idx: 2, Prefix: "-", Name: "h"},
			},
			expect: ValueOption{
				Option: options[0],
				Tok:    flagscanner.OptionToken{Idx: 2, Prefix: "-", Name: "h"},
			},
		},

		{
			name: "no early options",
			tokens: []flagscanner.Token{
				flagscanner.OptionToken{Idx: 0, Prefix: "-", Name: "x"},
				flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
			},
			expect: nil,
		},

		// The scanner will transform everything after a separator into
		// a positional argument so we should not parse --help
		{
			name: "with --help after a separator",
			tokens: []flagscanner.Token{
				flagscanner.OptionToken{Idx: 0, Prefix: "-", Name: "x"},
				flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
				flagscanner.OptionsArgumentsSeparatorToken{Idx: 2, Separator: "--"},
				flagscanner.PositionalArgumentToken{Idx: 3, Value: "--help"},
			},
			expect: nil,
		},

		// The scanner will transform everything after a separator into
		// a positional argument so we should not parse --help
		{
			name: "with -h after a separator",
			tokens: []flagscanner.Token{
				flagscanner.OptionToken{Idx: 0, Prefix: "-", Name: "x"},
				flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
				flagscanner.OptionsArgumentsSeparatorToken{Idx: 2, Separator: "--"},
				flagscanner.PositionalArgumentToken{Idx: 3, Value: "-h"},
			},
			expect: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value, found := earlyParse(options, tc.tokens)
			expectFound := tc.expect != nil
			assert.True(t, found == expectFound)
			assert.Equal(t, tc.expect, value)
		})
	}
}
