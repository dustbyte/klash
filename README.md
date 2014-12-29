klash
=====

### A new way to have arguments

*klash* exposes a new way of dealing with command line argument parsing in Golang using type reflection.

Indeed, parsing command line arguments is as easy as declaring and instanciate a Golang structure and hammering it with *klash*. It will result automatic association based on the provided structure fields for non-positional arguments, as well as a list remaining arguments.

Installation
------------
```shell
$ go get github.com/mota/klash
```

How to use
----------

First thing first, here is the prototype of the `Parse` function

```go
func Parse(params interface{}) ([]string, error)
```

Params must be a pointer to a structure type. What fields compose the structure is left to your discretion, provided that types are supported.

Suppose we want to greet people in a certain way

```go
package main

import (
	"fmt"
	"strings"

	"github.com/mota/klash"
)

type Arguments struct {
	Version  bool
	Greeting string
	Names    []string
}

func main() {
	args := Arguments{Greeting: "Hi"}
	leftover, err := klash.Parse(&args)

	if err == nil {
		if args.Version {
			fmt.Println("Greeting version 1.0")
			return
		}
		for _, name := range args.Names {
			fmt.Printf("%s %s!\n", args.Greeting, name)
		}
		if len(leftover) > 0 {
			fmt.Printf("Also %s to %s!\n",
				strings.ToLower(args.Greeting),
				strings.Join(leftover, " and "))
		}
	}
}
```

Which results to
```shell
$ ./main --greeting Hello --names Paul --names Ringo John George
Hello Paul!
Hello Ringo!
Also hello to John and George!
```

Simple as that!

Features
--------

### Supported
- Parse any of the given types: `int`, `uint`, `float32`, `float64`, `string`, `bool`
- Parse slices of supported types by repeating the argument name
- Optional stopped parsing when first positional argument is encountered (can be set with function `ParseArguments`, function `Parse` stops by default)

### In the future
- Automatic help and usage generator
- Subparsing via sub structures
- Golang Tag handling to support via a simple format
  - Help message
  - Alias names
  - Exclusion
  - Groups
  - Requirement
  - Coerce positional arguments
- More types (`time.Time`, other integer types ...)
