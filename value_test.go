//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/value_test.go
//

package flagparser

import (
	"testing"

	"github.com/bassosimone/flagscanner"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	// Just a random token for testing the Token method
	testToken := flagscanner.PositionalArgumentToken{
		Idx:   155,
		Value: "antani",
	}

	// Test case definition
	type testcase struct {
		name    string
		input   Value
		strings []string
		panics  bool
	}

	cases := []testcase{
		{
			name: "OptionTypeEarlyArgumentNone",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "help",
					Type:         OptionTypeEarlyArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"-help"},
			panics:  false,
		},

		{
			name: "OptionTypeGroupableArgumentNone",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "z",
					Type:         OptionTypeGroupableArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"-z"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentNone",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "verbose",
					Type:         OptionTypeStandaloneArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"--verbose"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentOptional",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "verbose",
					Type:         OptionTypeStandaloneArgumentOptional,
				},
				Value: "false",
			},
			strings: []string{"--verbose=false"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentRequired",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "file",
					Type:         OptionTypeStandaloneArgumentRequired,
				},
				Value: "/dev/null",
			},
			strings: []string{"--file", "/dev/null"},
			panics:  false,
		},

		{
			name: "OptionTypeGroupableArgumentRequired",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "o",
					Type:         OptionTypeGroupableArgumentRequired,
				},
				Value: "/dev/null",
			},
			strings: []string{"-o", "/dev/null"},
			panics:  false,
		},

		{
			name: "OptionType_invalid",
			input: ValueOption{
				Tok: testToken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "o",
					Type:         0, // invalid
				},
				Value: "/dev/null",
			},
			strings: []string{},
			panics:  true,
		},

		{
			name: "ValuePositionalArgument",
			input: ValuePositionalArgument{
				Tok:   testToken,
				Value: "/dev/null",
			},
			strings: []string{"/dev/null"},
			panics:  false,
		},

		{
			name: "ValueOptionsArgumentsSeparator",
			input: ValueOptionsArgumentsSeparator{
				Separator: "--",
				Tok:       testToken,
			},
			strings: []string{"--"},
			panics:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics {
				assert.Panics(t, func() {
					_ = tc.input.Strings()
				})
				return
			}
			gotStrings := tc.input.Strings()
			assert.Equal(t, tc.strings, gotStrings)
			assert.Equal(t, testToken, tc.input.Token())
		})
	}
}

func Test_sortValues(t *testing.T) {
	input := []Value{
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 2}, Value: "b"},
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 1}, Value: "a"},
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 3}, Value: "c"},
	}

	expected := []Value{
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 1}, Value: "a"},
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 2}, Value: "b"},
		ValuePositionalArgument{Tok: flagscanner.PositionalArgumentToken{Idx: 3}, Value: "c"},
	}

	sortValues(input)

	assert.Equal(t, expected, input)
}
