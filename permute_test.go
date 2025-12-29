//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/permute_test.go
//

package flagparser

import (
	"testing"

	"github.com/bassosimone/flagscanner"
	"github.com/stretchr/testify/assert"
)

func Test_maybePermute(t *testing.T) {

	// The command line arguments we are dealing with is the following:
	//
	//	-v testlist7.txt testlist111.txt --logs logs.jsonl -o output.txt testlist444.txt -- curl -o /dev/null

	options := []Value{
		ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: flagscanner.OptionToken{Idx: 1}},
		ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: flagscanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
		ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: flagscanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},
	}

	positionals := []Value{
		ValuePositionalArgument{Value: "testlist7.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 2}},
		ValuePositionalArgument{Value: "testlist111.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 3}},
		ValuePositionalArgument{Value: "testlist444.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 8}},
		ValueOptionsArgumentsSeparator{Separator: "--", Tok: flagscanner.OptionsArgumentsSeparatorToken{Idx: 9}},
		ValuePositionalArgument{Value: "curl", Tok: flagscanner.PositionalArgumentToken{Idx: 10}},
		ValuePositionalArgument{Value: "-o", Tok: flagscanner.PositionalArgumentToken{Idx: 11}},
		ValuePositionalArgument{Value: "/dev/null", Tok: flagscanner.PositionalArgumentToken{Idx: 12}},
	}

	// Define the test cases
	type testcase struct {
		name    string
		disable bool
		expect  []Value
	}
	cases := []testcase{
		{
			name:    "with permutation",
			disable: false,
			expect: []Value{
				// sorted options
				ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: flagscanner.OptionToken{Idx: 1}},
				ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: flagscanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
				ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: flagscanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},

				// sorted positional arguments
				ValuePositionalArgument{Value: "testlist7.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 2}},
				ValuePositionalArgument{Value: "testlist111.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 3}},
				ValuePositionalArgument{Value: "testlist444.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 8}},
				ValueOptionsArgumentsSeparator{Separator: "--", Tok: flagscanner.OptionsArgumentsSeparatorToken{Idx: 9}},
				ValuePositionalArgument{Value: "curl", Tok: flagscanner.PositionalArgumentToken{Idx: 10}},
				ValuePositionalArgument{Value: "-o", Tok: flagscanner.PositionalArgumentToken{Idx: 11}},
				ValuePositionalArgument{Value: "/dev/null", Tok: flagscanner.PositionalArgumentToken{Idx: 12}},
			},
		},

		{
			name:    "without permutation",
			disable: true,
			expect: []Value{
				ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: flagscanner.OptionToken{Idx: 1}},
				ValuePositionalArgument{Value: "testlist7.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 2}},
				ValuePositionalArgument{Value: "testlist111.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 3}},
				ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: flagscanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
				ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: flagscanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},
				ValuePositionalArgument{Value: "testlist444.txt", Tok: flagscanner.PositionalArgumentToken{Idx: 8}},
				ValueOptionsArgumentsSeparator{Separator: "--", Tok: flagscanner.OptionsArgumentsSeparatorToken{Idx: 9}},
				ValuePositionalArgument{Value: "curl", Tok: flagscanner.PositionalArgumentToken{Idx: 10}},
				ValuePositionalArgument{Value: "-o", Tok: flagscanner.PositionalArgumentToken{Idx: 11}},
				ValuePositionalArgument{Value: "/dev/null", Tok: flagscanner.PositionalArgumentToken{Idx: 12}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := permute(tc.disable, options, positionals)
			assert.Equal(t, tc.expect, got)
		})
	}
}
