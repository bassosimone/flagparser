//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/config_test.go
//

package flagparser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bassosimone/flagscanner"
	"github.com/stretchr/testify/assert"
)

func TestErrUnknownOption(t *testing.T) {
	err := ErrUnknownOption{
		Name:   "verbose",
		Prefix: "--",
		Token: flagscanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "verbose",
		},
	}
	expect := "unknown option: --verbose"
	assert.Equal(t, expect, err.Error())
}

func TestErrAmbiguousPrefix(t *testing.T) {
	err := ErrAmbiguousPrefix{
		Prefix: "-",
	}

	expect := `prefix "-" is used for both standalone and groupable options`
	assert.Equal(t, expect, err.Error())
}

func TestErrMultipleOptionsWithSameName(t *testing.T) {
	expect := `multiple options with "foo" name`

	opt1 := &Option{Name: "foo"}
	opt2 := &Option{Name: "foo"}
	err := ErrMultipleOptionsWithSameName{
		Name:    "foo",
		Options: []*Option{opt1, opt2},
	}

	assert.Equal(t, expect, err.Error())
}

func TestErrTooLongGroupableOptionName(t *testing.T) {
	opt := &Option{Name: "longname"}
	err := ErrTooLongGroupableOptionName{Option: opt}

	expect := "groupable option names should be a single byte, found: &{DefaultValue: Prefix: Name:longname Type:0}"
	assert.Equal(t, expect, err.Error())
}

func TestErrEmptyOptionName(t *testing.T) {
	opt := &Option{Name: ""}
	err := ErrEmptyOptionName{Option: opt}

	expect := "option name cannot be empty: &{DefaultValue: Prefix: Name: Type:0}"
	assert.Equal(t, expect, err.Error())
}

func TestErrEmptyOptionPrefix(t *testing.T) {
	opt := &Option{Prefix: ""}
	err := ErrEmptyOptionPrefix{Option: opt}

	expect := "option prefix cannot be empty: &{DefaultValue: Prefix: Name: Type:0}"
	assert.Equal(t, expect, err.Error())
}

func Test_config_disablePermute(t *testing.T) {
	cases := []bool{true, false}
	for _, tc := range cases {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			cfg := config{parser: &Parser{DisablePermute: tc}}
			assert.Equal(t, cfg.parser.DisablePermute, cfg.disablePermute())
		})
	}
}

func Test_config_findOption(t *testing.T) {
	// Create the option we would like to return to the caller
	option := Option{
		DefaultValue: "",
		Prefix:       "--",
		Name:         "verbose",
		Type:         OptionTypeStandaloneArgumentNone,
	}

	// Create a parser with a single option inside
	cfg := config{
		options: map[string]*Option{
			"verbose": &option,
		},
	}

	// Define the test cases
	type testcase struct {
		caseName         string
		tok              flagscanner.OptionToken
		optName          string
		kind             OptionType
		expectOp         *Option
		expectErrUnknown bool
	}
	cases := []testcase{
		{
			caseName: "successful find",
			tok: flagscanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindStandalone,
			expectOp:         &option,
			expectErrUnknown: false,
		},

		{
			caseName: "no such option",
			tok: flagscanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "file",
			kind:             optionKindStandalone,
			expectOp:         nil,
			expectErrUnknown: true,
		},

		{
			caseName: "the prefix does not match",
			tok: flagscanner.OptionToken{
				Idx:    4,
				Prefix: "/",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindStandalone,
			expectOp:         nil,
			expectErrUnknown: true,
		},

		{
			caseName: "the option kind does not match",
			tok: flagscanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindEarly,
			expectOp:         nil,
			expectErrUnknown: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			gotOption, err := cfg.findOption(tc.tok, tc.optName, tc.kind)

			switch {
			case tc.expectErrUnknown:
				var errval ErrUnknownOption
				good := errors.As(err, &errval)
				assert.True(t, good)
				assert.Equal(t, tc.optName, errval.Name)
				assert.Equal(t, tc.tok.Prefix, errval.Prefix)
				assert.Equal(t, tc.tok, errval.Token)

			case err != nil:
				t.Fatal(err)

			default:
				assert.Equal(t, tc.expectOp, gotOption)
			}
		})
	}
}

func Test_newConfig(t *testing.T) {
	// Define the structure of the test cases
	type testcase struct {
		caseName       string                // Name of the test case
		options        []*Option             // Options to be used in the parser
		expectErr      error                 // Expected error, if any
		expectPrefixes map[string]OptionType // Expected prefixes and their types
		expectOptions  map[string]*Option    // Expected options by name
	}

	// Define the test cases
	cases := []testcase{
		{
			caseName: "groupable option with multi-byte name",
			options: []*Option{
				{
					Name:   "longname",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectErr: ErrTooLongGroupableOptionName{
				Option: &Option{
					Name:   "longname",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "empty option name",
			options: []*Option{
				{
					Name:   "",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrEmptyOptionName{
				Option: &Option{
					Name:   "",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "empty option prefix",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrEmptyOptionPrefix{
				Option: &Option{
					Name:   "verbose",
					Prefix: "",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "multiple options with same name",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "verbose",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrMultipleOptionsWithSameName{
				Name: "verbose",
				Options: []*Option{
					{
						Name:   "verbose",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentNone,
					},
					{
						Name:   "verbose",
						Prefix: "-",
						Type:   OptionTypeStandaloneArgumentNone,
					},
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "ambiguous parsing prefixes",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectErr: ErrAmbiguousPrefix{
				Prefix: "-",
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "valid configuration",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
				{
					Name:   "help",
					Prefix: "--",
					Type:   OptionTypeEarlyArgumentNone,
				},
				{
					Name:   "h",
					Prefix: "-",
					Type:   OptionTypeEarlyArgumentNone,
				},
			},
			expectErr: nil,
			expectPrefixes: map[string]OptionType{
				"--": optionKindStandalone | optionKindEarly,
				"-":  optionKindGroupable | optionKindEarly,
			},
			expectOptions: map[string]*Option{
				"h": {
					Name:   "h",
					Prefix: "-",
					Type:   OptionTypeEarlyArgumentNone,
				},
				"help": {
					Name:   "help",
					Prefix: "--",
					Type:   OptionTypeEarlyArgumentNone,
				},
				"verbose": {
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				"v": {
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
		},
	}

	// Run through each test case
	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			// Create a parser with the provided options
			parser := &Parser{Options: tc.options}

			// Attempt to create a new config
			cfg, err := newConfig(parser)

			// Check for expected error
			assert.Equal(t, tc.expectErr, err)
			if err != nil {
				return
			}

			// Check the prefixes and options in the config
			assert.Equal(t, tc.expectPrefixes, cfg.prefixes)
			assert.Equal(t, tc.expectOptions, cfg.options)
		})
	}
}
