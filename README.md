klash
=====

### A new way to have arguments

*klash* exposes a new way of dealing with command line argument parsing in Golang using type reflection.

Indeed, parsing command line arguments is as easy as declaring and instanciate a Golang structure and hammering it with *klash*. It will result automatic association based on the provided structure fields for non-positional arguments, as well as a list of remaining arguments.

Installation
------------
```shell
$ go get github.com/mota/klash
```

How to use
----------

First thing first, here is the prototype of the `Parse` function

```go
func Parse(parameters interface{}, help string) []string
```

Params must be a pointer to a structure type. What fields compose the structure is left to your discretion, provided that types are supported or convertible.

Help is a description message of the application, it can be left to a void `string` value.

Suppose we want to greet people:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/mota/klash"
)

type Arguments struct {
	Version  bool     `klash-alias:"v" klash-help:"Print version and exit"`
	Greeting string   `klash-help:"Greeting expression"`
	Names    []string `klash-help:"List of people to be greeted"`
}

func main() {
	args := Arguments{Greeting: "Hi"} // Default values are set with Go syntactic constructions
	leftover := klash.Parse(&args, "Greet people")

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
```

Which results to
```shell
$ ./main --greeting Hello --names Paul --names Ringo John George
Hello Paul!
Hello Ringo!
Also hello to John and George!
```

Simple as that!

Documentation
-------------

Supported types are:
- `int`, `uint`
- `float32`, `float64`
- `string`
- `bool`

In addition to these types, any type *convertible* is taken into account.
What *convertible* means is that the type implements the following interface:

```go
type Convertible interface {
	FromString(stringval string) error
}
```

The `FromString` method receives a string representation of the value given as an argument. The concerned type is supposed to fill-up itself with the given data. If conversion fails (e.g the format is incomplete), the function shall return a descriptive error of what happened.

A basic example would be as follow:

```go
type Date struct {
	time.Time
}

func (d *Date) FromString(stringval string) error {
	val, err := time.Parse("2006-01-02", stringval)

	if err != nil {
		return err
	}

	d.Time = val
	return nil
}
```

Moreover, *klash* allows to store arguments into `slices` of any of the given types, even *convertible* ones.

#### Tags

As seen in the example above, Golang tags can be specified for each argument in order to add some information about the argument or its behavior.

At the moment, the tags are:

- `klash-help`: which allows to inform about the nature of the argument
- `klash-alias`: which allows to specify an alias name for the argument. Usually the alias name is shorter than the argument's name, but this is not required.

#### Help

Help messages as well as help switches (`-h` and `--help`) are handled automatically. Upon failure the error message is also displayed in the *standard error* stream.
The usage is generated automatically as well as the list of arguments and a representation of their types. What is left to the user are the general description message that is provided to `klash.Parse`, as well as arguments help messages as seen in the tags section.

#### Behavior

*klash* tries to use a behavior as wide as possible, but it might be possible that this acting doesn't fit what your're used to do.

First of all, as does the `flag` Go package, by default *klash* stops making associations when the first positional argument is encountered.  The following example `./prog --delete --posterior-to 2014-07-31 *~ --version` will take into account the argument `--version` as a positional argument. If you wish to flip the *version* flag, do it before the `*~` argument.

Secondly, provided they have the length of their name/alias one-character long, boolean argument can be packed into one single non-positional argument on call. That is, if you have two boolean arguments `Debug` and `Version` which are aliased to `d` and `v`, you can call them with the more terse `-dv`. Or even `-vd`, it's up to the user.

Afterwards, the number of dashes preceeding a positional argument name as well as its case are not important.

And finally (I'm running out of synonyms), you can interchangeably use the space-separated form as well as the assignment form to store values into arguments. For example: `./main --names=John --names George` is perfectly valid.

Having a finer control
----------------------

By default, *klash* makes choices for you in order to provide a simple a clear interface. These default choices are trusted to follow some sort of common-sense. However it may not reflect the desired behavior of everyone.

That is why *klash* also provides a little more complex but more powerful interface:

```go
func ParseArguments(arguments []string, parameters interface{}, stop bool) ([]string, error)
```

- `arguments` is the list of arguments given to the program, usually `os.Args[1:]`
- `parameters`, as well as the `Parse` interface, this parameter is a pointer to the structure that defines your arguments
- `stop` which affects the *stop* behavior. That is, if `stop` is set to false, the parser will continue to take account of non-positional arguments until there is no more argument to examine

This function returns two values:
- A slice of strings `[]string` which, as well as the `Parse` function, contains the positional arguments
- A `error` value that describes textually the error produced. If the error value is the same as `klash.HelpError` therefore a help switch has been activated. Albeit in one hand this does mean that no sensitive error has been detected, on the other hand it doesn't mean the argument parsing is terminated nor it wouldn't have produced legitimate errors. Take for granted this is an erroneous behavior in any case.

#### Subparsing
In future versions, *klash* will provide a means to handle subparsing (a.k.a subcommands) automatically within the arguments structure. However, since a structure is is sometimes too monolithic, this rigidity can be inconvenient for subparsing.

Nevertheless, the `ParseArguments` functions allows to extend parsing at runtime and eases subcommand dispatching.
How does it work? Pretty simple.

Suppose you are to develop a simple http client over a REST API, there will be commands such as `push`, `pull` or even `update`. Let also say you want to set the verbose mode globally.

```go
type GlobalArgs struct {
	Verbose bool `klash-alias:"v"`
}
args := GlobalArgs{} // bool default is to false
leftover, _ := klash.ParseArguments(os.Args[1:], &args, true) // stop is set to true
```

Now we just have to check the first argument of the leftover slice, which will contain the name of the subcommand to be executed.

```go
cmd := leftover[0]
switch cmd {
// ...
case "pull":
	type PullArgs struct {
		Tags []string
	}
	pa := PullArgs{}
	resources, _ := klash.ParseArguments(leftover[1:], &pa, true)
	engine.pull(args.Verbose, pa.Tags, resources)
}
```

Provided a real-case implementing this kind of behavior, a simple call as follow should do the work:
```
$ ./client -v pull --tags "debian" --tags "openbsd"
```

Et voilà !
