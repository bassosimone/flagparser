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
