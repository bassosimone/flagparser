# Golang Low-Level CLI Flags Parser

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/flagparser)](https://pkg.go.dev/github.com/bassosimone/flagparser) [![Build Status](https://github.com/bassosimone/flagparser/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/flagparser/actions) [![codecov](https://codecov.io/gh/bassosimone/flagparser/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/flagparser)

The `flagparser` Go package contains a low-level command-line
arguments parser. It is a building block that enables building
higher-level command-line-flags parsers.

For example:

```Go
import (
	"log"
	"os"

	"github.com/bassosimone/flagparser"
)

// Construct a parser recognizing GNU style options.
p := &flagparser.NewParser()
p.SetMinMaxPositionalArguments(1, 1)                    // exactly one positional
p.AddLongOptionWithArgumentOptional("compress", "gzip") // --compress[=gzip|bzip2|...]
p.AddEarlyOption('h', "help")                           // -h, --help
p.AddOptionWithArgumentRequired('o', "output")          // -o, --output FILE
p.AddOptionWithArgumentNone('v', "verbose")             // -v, --verbose

// Parse the command line
values, err := p.Parse(os.Args[1:])
if err != nil {
	log.Fatal(err)
}
```

The above example configures GNU style options but we support a
wide variety of command-line-flags styles including Go, dig, Windows,
and traditional Unix. See [example_test.go](example_test.go).

## Installation

To add this package as a dependency to your module:

```sh
go get github.com/bassosimone/flagparser
```

## Development

To run the tests:
```sh
go test -v .
```

To measure test coverage:
```sh
go test -v -cover .
```

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```

## History

Adapted from [bassosimone/clip](https://github.com/bassosimone/clip/tree/v0.8.0).
