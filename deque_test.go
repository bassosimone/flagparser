//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/deque_test.go
//

package flagparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_deque(t *testing.T) {
	// Start with a deque containing two elements
	original := []Value{
		ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Value: "FILE"},
		ValuePositionalArgument{Value: "http://www.google.com/"},
	}
	input := deque[Value]{values: original}

	// Extract from the deque like we're going to do when parsing
	var output deque[Value]
	for !input.Empty() {
		value, good := input.Front()
		if !good {
			t.Fatal("expected to be able to extract an element")
		}
		input.PopFront()
		output.PushBack(value)
	}

	// Compare the results
	assert.Equal(t, original, output.values)
}
