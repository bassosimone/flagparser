//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/doc.go
//

/*
Package flagparser implements a flexible command line parser.

[NewParser] configures GNU-style defaults: short options use `-`, long options
use `--`, the options-arguments separator is `--`, and option permutation is
enabled. You can override any of these defaults to parse non-GNU command lines.

To parse arguments, you need to:

 1. Create a [*Parser] instance (typically using the [NewParser] factory).

 2. Initialize its options (e.g., with [*Parser.AddOptionWithArgumentNone] or by
    appending [Option] values to the [*Parser] instance.).

 3. Optionally, adjust the separator, permutation, and prefixes to match the
    desired command-line convention.

 4. Invoke [*Parser.Parse] passing it `os.Args[1:]`.

The [*Parser.Parse] method returns a slice of [Value].

# Options-Arguments Separator

The [*Parser] can be configured to define a separator after which any
command line token is treated as a positional argument, regardless
of its prefix. The GNU getopt implementation and the Go standard library do
this using the `--` separator. [NewParser] configures `--` as separator.

# Permutation

By default, the parser permutes options ahead of positional arguments,
matching the GNU getopt behavior. You can disable permutation (see
the [*Parser.DisablePermute] knob) to preserve the original order, which
can be useful when a subcommand expects its own flags.

# Option Types

Each [Option] has its own [OptionType], which is one of these values:

 1. [OptionTypeEarlyArgumentNone]: options processed before the actual command-line
    parsing to detect flags (e.g., `--help`) that should always cause specific
    actions (e.g., printing the help message on the stdout), regardless
    of the correctness of the rest of the command line. These options cannot
    receive arguments since they are processed ahead of the parsing.

 2. [OptionTypeStandaloneArgumentNone]: options that cannot be grouped
    and that require no arguments (e.g., `--verbose`).

 3. [OptionTypeStandaloneArgumentRequired]: options that cannot be grouped
    and that require an argument. The argument can be provided in a subsequent
    token (e.g., `--file FILE`) or after the `=` byte (`--file=FILE`).

 4. [OptionTypeStandaloneArgumentOptional]: options that cannot be grouped
    and take an optional argument. The argument must be provided after
    the `=` byte (e.g., `--deepscan=true`, `--deepscan=false`). Omitting the
    value (e.g., `--deepscan`) causes the default value to be used.

 5. [OptionTypeGroupableArgumentNone]: single-letter options that can be
    grouped together (e.g., `-xz` as a shortcut for `-x -z`).

 6. [OptionTypeGroupableArgumentRequired]: like the previous section but an
    argument must be specified, either as a subsequent token (e.g.,
    `-xzf FILE`) or directly after the option (`-xzfFILE`) -- note that
    even though the latter may be confusing it is a GNU extension.

# Option Prefixes

Each [Option] can define its own parsing prefix. Generally, it is
advisable to use uniform prefixes for all options. For example, following
the GNU convention, one should use the `-` prefix for groupable options
and the `--` prefix for standalone options. This is the behavior that you
get if you use [*Parser.AddOptionWithArgumentNone] and similar functions
to initialize a parser. However, you can also use non-GNU conventions, such
as standalone options prefixed by `-`, thus emulating the Go flag package
option parsing style. To this end, manually create the [Option].

This package also supports using distinct prefixes for distinct
options of the same type. For example, both `+short` and `--verbose`
could be standalone options. The only restriction, enforced by
the [*Parser], is that you cannot use the same prefix for groupable
and standalone options. That is, if `-` is used for groupable
options it cannot be used for standalone options as well.

The early options are an exception to this rule, since they
are not really parsed, rather just pattern matched against the
argv provided by the programmer. Therefore, it is possible to
have `-a` and `-b` as groupable options and `-h` for help,
provided that you declare `-h` as an early option. In other
words, the prefixes assigned to early options do not have
an impact on the single-prefix restriction.

# Parsed Values

 1. [ValueOption]: contains a parsed [*Option].

 2. [ValuePositionalArgument]: contains a positional argument.

 3. [ValueOptionsArgumentsSeparator]: contains the separator
    between the options and the arguments (usually `--`).

# Example

Consider the following command line arguments:

	-sv --output /dev/null -- https://example.com/

Assume you define these options:

	Option{Name:s Prefix:- Type:OptionTypeGroupableArgumentNone}
	Option{Name:v Prefix:- Type:OptionTypeGroupableArgumentNone}
	Option{Name:output Prefix:-- Type:OptionTypeStandaloneArgumentRequired}

Assume you use `--` as the options-arguments separator.

Then, the parser will return:

	ValueOption{Token:scanner.TokenOption{Prefix:-} Name:s}
	ValueOption{Token:scanner.TokenOption{Prefix:-} Name:v}
	ValueOption{Token:scanner.TokenOption{Prefix:--} Name:output Value:/dev/null}
	ValueOptionsArgumentsSeparator{}
	ValuePositionalArgument{Value:https://example.com/}

See the package examples for more examples.
*/
package flagparser
