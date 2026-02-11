//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/early.go
//

package flagparser

import "github.com/bassosimone/flagscanner"

// earlyParse parses the early options. That is, the options that should
// be recognized immediately even when the rest of the command line is
// wrong. The `--help` option is the most typical early option we handle.
//
// When disablePermute is true, we stop scanning as soon as we encounter
// a positional argument, mirroring the normal parsing behavior where a
// positional stops option recognition.
func earlyParse(options []*Option, tokens []flagscanner.Token, disablePermute bool) (Value, bool) {
	// 1. process each token and only consider the option tokens
	for _, tok := range tokens {
		switch tok := tok.(type) {
		case flagscanner.OptionToken:
			// 2. process each option
			for _, option := range options {
				if option.Type.isEarly() && tok.Prefix == option.Prefix && tok.Name == option.Name {

					// We have found the early option, return it
					eopt := ValueOption{
						Option: option,
						Tok:    tok,
						Value:  "",
					}
					return eopt, true

				}
			}

		case flagscanner.PositionalArgumentToken:
			if disablePermute {
				return nil, false
			}
		}
	}
	return nil, false
}
