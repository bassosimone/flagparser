//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/permute.go
//

package flagparser

func permute(disablePermute bool, options, positionals []Value) []Value {
	// Determine what to do depending on the configuration
	switch {

	// When permutation is disabled, restore the original token order
	case disablePermute:
		output := make([]Value, 0, 1+len(options)+len(positionals))
		output = append(output, options...)
		output = append(output, positionals...)
		sortValues(output)
		return output

	// Otherwise, merge options together and sort options and arguments independently
	default:
		sortValues(options)
		sortValues(positionals)
		output := make([]Value, 0, 1+len(options)+len(positionals))
		output = append(output, options...)
		output = append(output, positionals...)
		return output
	}
}
