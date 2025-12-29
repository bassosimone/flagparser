//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/optiontype_test.go
//

package flagparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionType(t *testing.T) {
	type testcase struct {
		name         string
		input        OptionType
		isEarly      bool
		isStandalone bool
		isGroupable  bool
	}

	cases := []testcase{
		{
			name:    "OptionTypeEarlyArgumentNone",
			input:   OptionTypeEarlyArgumentNone,
			isEarly: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentNone",
			input:        OptionTypeStandaloneArgumentNone,
			isStandalone: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentRequired",
			input:        OptionTypeStandaloneArgumentRequired,
			isStandalone: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentOptional",
			input:        OptionTypeStandaloneArgumentOptional,
			isStandalone: true,
		},

		{
			name:        "OptionTypeGroupableArgumentNone",
			input:       OptionTypeGroupableArgumentNone,
			isGroupable: true,
		},

		{
			name:        "OptionTypeGroupableArgumentRequired",
			input:       OptionTypeGroupableArgumentRequired,
			isGroupable: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isEarly, tc.input.isEarly())
			assert.Equal(t, tc.isStandalone, tc.input.isStandalone())
			assert.Equal(t, tc.isGroupable, tc.input.isGroupable())
		})
	}
}

func Test_NewOptionWithArgumentNone(t *testing.T) {
	t.Run("short only", func(t *testing.T) {
		options := NewOptionWithArgumentNone('v', "")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "v",
				Type:   OptionTypeGroupableArgumentNone,
			}, options[0])
		}
	})

	t.Run("long only", func(t *testing.T) {
		options := NewOptionWithArgumentNone(0, "verbose")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "verbose",
				Type:   OptionTypeStandaloneArgumentNone,
			}, options[0])
		}
	})

	t.Run("short and long", func(t *testing.T) {
		options := NewOptionWithArgumentNone('v', "verbose")
		if assert.Len(t, options, 2) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "v",
				Type:   OptionTypeGroupableArgumentNone,
			}, options[0])
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "verbose",
				Type:   OptionTypeStandaloneArgumentNone,
			}, options[1])
		}
	})

	t.Run("no options", func(t *testing.T) {
		options := NewOptionWithArgumentNone(0, "")
		assert.Nil(t, options)
	})
}

func Test_NewEarlyOption(t *testing.T) {
	t.Run("short only", func(t *testing.T) {
		options := NewEarlyOption('h', "")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "h",
				Type:   OptionTypeEarlyArgumentNone,
			}, options[0])
		}
	})

	t.Run("long only", func(t *testing.T) {
		options := NewEarlyOption(0, "help")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "help",
				Type:   OptionTypeEarlyArgumentNone,
			}, options[0])
		}
	})

	t.Run("short and long", func(t *testing.T) {
		options := NewEarlyOption('h', "help")
		if assert.Len(t, options, 2) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "h",
				Type:   OptionTypeEarlyArgumentNone,
			}, options[0])
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "help",
				Type:   OptionTypeEarlyArgumentNone,
			}, options[1])
		}
	})

	t.Run("no options", func(t *testing.T) {
		options := NewEarlyOption(0, "")
		assert.Nil(t, options)
	})
}

func Test_NewOptionWithArgumentRequired(t *testing.T) {
	t.Run("short only", func(t *testing.T) {
		options := NewOptionWithArgumentRequired('o', "")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "o",
				Type:   OptionTypeGroupableArgumentRequired,
			}, options[0])
		}
	})

	t.Run("long only", func(t *testing.T) {
		options := NewOptionWithArgumentRequired(0, "output")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "output",
				Type:   OptionTypeStandaloneArgumentRequired,
			}, options[0])
		}
	})

	t.Run("short and long", func(t *testing.T) {
		options := NewOptionWithArgumentRequired('o', "output")
		if assert.Len(t, options, 2) {
			assert.Equal(t, &Option{
				Prefix: "-",
				Name:   "o",
				Type:   OptionTypeGroupableArgumentRequired,
			}, options[0])
			assert.Equal(t, &Option{
				Prefix: "--",
				Name:   "output",
				Type:   OptionTypeStandaloneArgumentRequired,
			}, options[1])
		}
	})

	t.Run("no options", func(t *testing.T) {
		options := NewOptionWithArgumentRequired(0, "")
		assert.Nil(t, options)
	})
}

func Test_NewLongOptionWithArgumentOptional(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		assert.Nil(t, NewLongOptionWithArgumentOptional("", "gzip"))
	})

	t.Run("long only", func(t *testing.T) {
		options := NewLongOptionWithArgumentOptional("compress", "gzip")
		if assert.Len(t, options, 1) {
			assert.Equal(t, &Option{
				DefaultValue: "gzip",
				Prefix:       "--",
				Name:         "compress",
				Type:         OptionTypeStandaloneArgumentOptional,
			}, options[0])
		}
	})
}
