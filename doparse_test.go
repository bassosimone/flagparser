//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/doparse_test.go
//

package flagparser

import (
	"errors"
	"testing"

	"github.com/bassosimone/flagscanner"
	"github.com/stretchr/testify/assert"
)

func TestErrOptionRequiresNoArgument(t *testing.T) {
	err := ErrOptionRequiresNoArgument{
		Option: &Option{
			DefaultValue: "",
			Prefix:       "--",
			Name:         "verbose",
			Type:         OptionTypeStandaloneArgumentNone,
		},
		Token: flagscanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "verbose",
		},
	}

	expect := "option requires no argument: --verbose"
	assert.Equal(t, expect, err.Error())
}

func TestErrOptionRequiresArgument(t *testing.T) {
	err := ErrOptionRequiresArgument{
		Option: &Option{
			DefaultValue: "",
			Prefix:       "--",
			Name:         "file",
			Type:         OptionTypeStandaloneArgumentRequired,
		},
		Token: flagscanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "file",
		},
	}

	expect := "option requires an argument: --file"
	assert.Equal(t, expect, err.Error())
}

func newTestDoParseConfig() *config {
	return &config{
		parser: &Parser{
			DisablePermute:            false,
			OptionsArgumentsSeparator: "--",
		},
		prefixes: map[string]OptionType{
			"--": optionKindStandalone,
			"-":  optionKindGroupable,
		},
		options: map[string]*Option{
			"file": {
				Prefix: "--",
				Name:   "file",
				Type:   OptionTypeStandaloneArgumentRequired,
			},
			"http": {
				DefaultValue: "1.1",
				Prefix:       "--",
				Name:         "http",
				Type:         OptionTypeStandaloneArgumentOptional,
			},
			"verbose": {
				Prefix: "--",
				Name:   "verbose",
				Type:   OptionTypeStandaloneArgumentNone,
			},
			"x": {
				Prefix: "-",
				Name:   "x",
				Type:   OptionTypeGroupableArgumentRequired,
			},
			"z": {
				Prefix: "-",
				Name:   "z",
				Type:   OptionTypeGroupableArgumentNone,
			},
		},
	}
}

func parseTokens(cfg *config, tokens []flagscanner.Token) ([]string, []string, error) {
	input := &deque[flagscanner.Token]{values: tokens}
	var options deque[Value]
	var positionals deque[Value]
	err := doParse(cfg, input, &options, &positionals)
	return flattenValues(options.values), flattenValues(positionals.values), err
}

func flattenValues(values []Value) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, value.Strings()...)
	}
	return out
}

func Test_doParse(t *testing.T) {
	t.Run("permute keeps parsing options after a positional", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.parser.DisablePermute = false
		opts, pos, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
			flagscanner.OptionToken{Idx: 2, Prefix: "--", Name: "verbose"},
			flagscanner.OptionToken{Idx: 3, Prefix: "-", Name: "x"},
			flagscanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"--verbose", "-x", "file2.txt"}, opts)
		assert.Equal(t, []string{"file1.txt"}, pos)
	})

	t.Run("disable permute turns later options into positionals", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.parser.DisablePermute = true
		opts, pos, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.PositionalArgumentToken{Idx: 1, Value: "file1.txt"},
			flagscanner.OptionToken{Idx: 2, Prefix: "--", Name: "verbose"},
		})
		assert.NoError(t, err)
		assert.Empty(t, opts)
		assert.Equal(t, []string{"file1.txt", "--verbose"}, pos)
	})

	t.Run("separator turns later options into positionals", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.parser.DisablePermute = false
		opts, pos, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "verbose"},
			flagscanner.OptionsArgumentsSeparatorToken{Idx: 2, Separator: "--"},
			flagscanner.OptionToken{Idx: 3, Prefix: "--", Name: "file"},
			flagscanner.PositionalArgumentToken{Idx: 4, Value: "file1.txt"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"--verbose"}, opts)
		assert.Equal(t, []string{"--", "--file", "file1.txt"}, pos)
	})

	t.Run("groupable and standalone option arguments", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.parser.DisablePermute = false
		opts, pos, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
			flagscanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
			flagscanner.OptionToken{Idx: 3, Prefix: "--", Name: "file"},
			flagscanner.PositionalArgumentToken{Idx: 4, Value: "/dev/null"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"-z", "-x", "file1.txt", "--file", "/dev/null"}, opts)
		assert.Empty(t, pos)
	})

	t.Run("optional argument default and explicit value", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.parser.DisablePermute = false
		opts, pos, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "http"},
			flagscanner.OptionToken{Idx: 2, Prefix: "--", Name: "http=2.0"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"--http=1.1", "--http=2.0"}, opts)
		assert.Empty(t, pos)
	})
}

func Test_doParse_errors(t *testing.T) {
	cfg := newTestDoParseConfig()
	cfg.parser.DisablePermute = false

	t.Run("unknown option", func(t *testing.T) {
		_, _, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "nope"},
		})
		var errvalue ErrUnknownOption
		assert.True(t, errors.As(err, &errvalue))
	})

	t.Run("missing required argument", func(t *testing.T) {
		_, _, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "file"},
		})
		var errvalue ErrOptionRequiresArgument
		assert.True(t, errors.As(err, &errvalue))
	})

	t.Run("argument provided to no-arg option", func(t *testing.T) {
		_, _, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "verbose=true"},
		})
		var errvalue ErrOptionRequiresNoArgument
		assert.True(t, errors.As(err, &errvalue))
	})

	t.Run("groupable missing required argument", func(t *testing.T) {
		_, _, err := parseTokens(cfg, []flagscanner.Token{
			flagscanner.OptionToken{Idx: 1, Prefix: "-", Name: "x"},
		})
		var errvalue ErrOptionRequiresArgument
		assert.True(t, errors.As(err, &errvalue))
	})
}

func Test_doParse_panics(t *testing.T) {
	t.Run("unhandled standalone option type", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.options["__panic"] = &Option{
			Prefix: "--",
			Name:   "__panic",
			Type:   optionKindStandalone,
		}
		assert.Panics(t, func() {
			_, _, _ = parseTokens(cfg, []flagscanner.Token{
				flagscanner.OptionToken{Idx: 1, Prefix: "--", Name: "__panic"},
			})
		})
	})

	t.Run("unhandled groupable option type", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		cfg.options["_"] = &Option{
			Prefix: "-",
			Name:   "_",
			Type:   optionKindGroupable,
		}
		assert.Panics(t, func() {
			_, _, _ = parseTokens(cfg, []flagscanner.Token{
				flagscanner.OptionToken{Idx: 1, Prefix: "-", Name: "_"},
			})
		})
	})

	t.Run("unbound prefix", func(t *testing.T) {
		cfg := newTestDoParseConfig()
		assert.Panics(t, func() {
			_, _, _ = parseTokens(cfg, []flagscanner.Token{
				flagscanner.OptionToken{Idx: 1, Prefix: "/", Name: "help"},
			})
		})
	})
}
