//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/bassosimone/clip/blob/v0.8.0/pkg/nparser/parser_test.go
//

package flagparser_test

import (
	"fmt"
	"log"
	"math"

	"github.com/bassosimone/flagparser"
	"github.com/bassosimone/runtimex"
)

// Successful parsing of curl-like invocation with short options where
// a required argument is provided as a separate token.
func Example_curlParsingSuccessShortWithSpace() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse; the `-o` argument is a separate token.
	argv := []string{"curl", "https://www.example.com/", "-fsSLo", "index.html"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [-f]
	// [-s]
	// [-S]
	// [-L]
	// [-o index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with short options where a required
// argument is glued to the last short flag (GNU extension).
func Example_curlParsingSuccessShortWithNoSpace() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse; the `-o` argument is glued to the flag.
	argv := []string{"curl", "https://www.example.com/", "-fsSLoindex.html"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [-f]
	// [-s]
	// [-S]
	// [-L]
	// [-o index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with long options where
// a required argument is provided as a separate token.
func Example_curlParsingSuccessLongWithSpace() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse; the `--output` argument is separate.
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output",
		"index.html",
	}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--fail]
	// [--silent]
	// [--show-error]
	// [--location]
	// [--output index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with long options where
// a required argument is provided after an '=' sign.
func Example_curlParsingSuccessLongWithEqual() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse; the `--output` argument uses `=`.
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
	}

	// Parse the options; the early option wins over other errors.
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--fail]
	// [--silent]
	// [--show-error]
	// [--location]
	// [--output index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with `-h` acting as an "early"
// option that short-circuits parsing regardless of other errors.
func Example_curlParsingSuccessWithEarlyHelpShort() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddEarlyOption('h', "help")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--nonexistent-option", // should cause failure
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
		"-h", // but we have `-h` here
	}

	// Parse the options; the early option wins over the nonexistent option
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [-h]
}

// Successful parsing of curl-like invocation with `--help` acting as an "early"
// option that short-circuits parsing regardless of other errors.
func Example_curlParsingSuccessWithEarlyHelpLong() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddEarlyOption('h', "help")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--nonexistent-option", // should cause failure
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
		"--help", // but we have `--help` here
	}

	// Parse the options; the early option wins over the nonexistent option
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--help]
}

// Failing parsing of curl-like invocation with too few positionals.
func Example_curlParsingFailureTooFewPositionalArguments() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse
	//
	// Note: we're not providing a URL.
	argv := []string{
		"curl",
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
	}

	// Parse the options; this is where min/max positionals are enforced.
	values, err := parser.Parse(argv[1:])
	runtimex.Assert(len(values) <= 0 && err != nil)

	// Print the error value
	fmt.Printf("%s\n", err.Error())

	// Output:
	// too few positional arguments: expected at least 1, got 0
}

// Failing parsing of curl-like invocation with too many positionals.
func Example_curlParsingFailureTooManyPositionalArguments() {
	// Define a parser accepting curl-like command line options.
	//
	// Note: the default is to expect zero positionals.
	parser := flagparser.NewParser()
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
	}

	// Parse the options; this is where min/max positionals are enforced.
	values, err := parser.Parse(argv[1:])
	runtimex.Assert(len(values) <= 0 && err != nil)

	// Print the error value
	fmt.Printf("%s\n", err.Error())

	// Output:
	// too many positional arguments: expected at most 0, got 1
}

// Failing parsing of curl-like invocation with an unknown option.
func Example_curlParsingFailureWithInvalidOption() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--nonexistent-option", // will cause failure
		"--fail",
		"--silent",
		"--show-error",
		"--location",
		"--output=index.html",
	}

	// Parse the options; this is where `--nonexistent-option` causes a failure
	values, err := parser.Parse(argv[1:])
	runtimex.Assert(len(values) <= 0 && err != nil)

	// Print the error value
	fmt.Printf("%s\n", err.Error())

	// Output:
	// unknown option: --nonexistent-option
}

// Successful parsing of git-submodule-foreach-like invocation with the
// options-positionals separator and option permutation enabled.
func Example_gitSubmoduleForeachWithSeparatorWithReordering() {
	// Define a parser accepting git-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(0, math.MaxInt)
	parser.AddOptionWithArgumentNone('r', "recursive")

	// Define the argument vector to parse
	argv := []string{"git", "submodule", "foreach", "--recursive", "--", "git", "status", "-v"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: reordering probably does not give us the desired output
	// unless the `--recursive` flag could be applied to `git`.
	//
	// You typically do not parse a command with subcommands directly
	// using this parser but this example is here to show what happens
	// and build a mental model of how it works.
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--recursive]
	// [submodule]
	// [foreach]
	// [--]
	// [git]
	// [status]
	// [-v]
}

// Successful parsing of git-submodule-foreach-like invocation with the
// options-positionals separator and option permutation disabled.
func Example_gitSubmoduleForeachWithSeparatorWithoutReordering() {
	// Define a parser accepting git-like command line options.
	parser := flagparser.NewParser()
	parser.DisablePermute = true
	parser.SetMinMaxPositionalArguments(0, math.MaxInt)
	parser.AddOptionWithArgumentNone('r', "recursive")

	// Define the argument vector to parse
	argv := []string{"git", "submodule", "foreach", "--recursive", "--", "git", "status", "-v"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note how this case gives us the correct positional result.
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [submodule]
	// [foreach]
	// [--recursive]
	// [--]
	// [git]
	// [status]
	// [-v]
}

// Failing parsing of git-submodule-foreach-like invocation without an explicit
// separator in the args and with option permutation enabled.
func Example_gitSubmoduleForeachWithoutSeparatorWithReordering() {
	// Define a parser accepting git-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(0, math.MaxInt)
	parser.AddOptionWithArgumentNone('r', "recursive")

	// Define the argument vector to parse
	argv := []string{"git", "submodule", "foreach", "--recursive", "git", "status", "-v"}

	// Parse the options
	//
	// Note: reordering probably does not give us the desired output
	// and specifically the `-v` is stolen from `git status` and causes
	// a parsing error, which is definitely not what we want.
	//
	// You typically do not parse a command with subcommands directly
	// using this parser but this example is here to show what happens
	// and build a mental model of how it works.
	values, err := parser.Parse(argv[1:])
	runtimex.Assert(len(values) <= 0 && err != nil)

	// Print the error value
	fmt.Printf("%s\n", err.Error())

	// Output:
	// unknown option: -v
}

// Successful parsing of git-submodule-foreach-like invocation without an explicit
// separator in the args and with option permutation disabled.
func Example_gitSubmoduleForeachWithoutSeparatorWithoutReordering() {
	// Define a parser accepting git-like command line options.
	parser := flagparser.NewParser()
	parser.DisablePermute = true
	parser.SetMinMaxPositionalArguments(0, math.MaxInt)
	parser.AddOptionWithArgumentNone('r', "recursive")

	// Define the argument vector to parse
	argv := []string{"git", "submodule", "foreach", "--recursive", "git", "status", "-v"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note how this case gives us the correct positional result.
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [submodule]
	// [foreach]
	// [--recursive]
	// [git]
	// [status]
	// [-v]
}

// Successful parsing of dig-like invocation mixing GNU-style `-p` with `+short`.
func Example_digParsingSuccessWithMixedPrefixes() {
	// Define a parser with mixed prefixes and a groupable required argument.
	parser := &flagparser.Parser{
		MinPositionalArguments: 1,
		MaxPositionalArguments: 4,
		Options: []*flagparser.Option{
			{
				Name:   "p",
				Prefix: "-",
				Type:   flagparser.OptionTypeGroupableArgumentRequired,
			},
			{
				Name:   "short",
				Prefix: "+",
				Type:   flagparser.OptionTypeStandaloneArgumentNone,
			},
			{
				DefaultValue: "1024",
				Name:         "bufsize",
				Prefix:       "+",
				Type:         flagparser.OptionTypeStandaloneArgumentOptional,
			},
		},
	}

	// Define the argument vector to parse; `+bufsize` uses the default value.
	argv := []string{"dig", "@8.8.8.8", "-p53", "IN", "+short", "+bufsize", "A", "example.com"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [-p 53]
	// [+short]
	// [+bufsize=1024]
	// [@8.8.8.8]
	// [IN]
	// [A]
	// [example.com]
}

// Successful parsing of curl-like invocation with long options where
// the optional option argument is not provided.
func Example_curlParsingSuccessLongWithOptionalValueNotPresent() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddLongOptionWithArgumentOptional("fail", "true")
	parser.AddOptionWithArgumentRequired('o', "output")

	// Define the argument vector to parse; the `--fail` argument uses the default value.
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--fail",
		"--output=index.html",
	}

	// Parse the options; the `--fail` gets assigned a default value.
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--fail=true]
	// [--output index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with long options where
// the optional option argument is provided after the `=` sign.
func Example_curlParsingSuccessLongWithOptionalValuePresent() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddLongOptionWithArgumentOptional("fail", "true")
	parser.AddOptionWithArgumentRequired('o', "output")

	// Define the argument vector to parse; the `--fail` argument uses an explicit value.
	argv := []string{
		"curl",
		"https://www.example.com/",
		"--fail=false",
		"--output=index.html",
	}

	// Parse the options; the `--fail` gets assigned a default value.
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [--fail=false]
	// [--output index.html]
	// [https://www.example.com/]
}

// Successful parsing of curl-like invocation with short options where
// a required argument value is the `-` separate token (commonly used to
// indicate that we want to use the standard output).
func Example_curlParsingSuccessShortWithSpaceAndDashValue() {
	// Define a parser accepting curl-like command line options.
	parser := flagparser.NewParser()
	parser.SetMinMaxPositionalArguments(1, math.MaxInt)
	parser.AddOptionWithArgumentNone('f', "fail")
	parser.AddOptionWithArgumentNone('L', "location")
	parser.AddOptionWithArgumentRequired('o', "output")
	parser.AddOptionWithArgumentNone('S', "show-error")
	parser.AddOptionWithArgumentNone('s', "silent")

	// Define the argument vector to parse; the `-o` argument is a separate token.
	argv := []string{"curl", "https://www.example.com/", "-fsSLo", "-"}

	// Parse the options
	values, err := parser.Parse(argv[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed values to stdout
	//
	// Note: we have reordering by default so options are sorted before
	// positional arguments (respecting their relative order)
	for _, value := range values {
		fmt.Printf("%+v\n", value.Strings())
	}

	// Output:
	// [-f]
	// [-s]
	// [-S]
	// [-L]
	// [-o -]
	// [https://www.example.com/]
}
