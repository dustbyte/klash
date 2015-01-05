package klash

import (
	"fmt"
	"os"
)

func ParseArguments(name string,
	arguments []string,
	parameters interface{},
	help string,
	stop bool) ([]string, error) {

	params, err := NewParams(parameters)
	if err != nil {
		return nil, err
	}

	parser := NewArgumentParser(params, arguments, stop)
	err = parser.Parse()

	if err != nil {
		return nil, err
	}

	return parser.OutArgs, nil
}

func HammerArguments(name string,
	arguments []string,
	parameters interface{},
	help string,
	stop bool) []string {

	params, err := NewParams(parameters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	parser := NewArgumentParser(params, arguments, stop)
	err = parser.Parse()

	if err != nil {
		status := 0

		fmt.Fprint(os.Stderr, GenerateHelp(name, help, params))

		if err != HelpError {
			status = 2
			fmt.Fprintf(os.Stderr, "\n%s: error: %s\n", name, err)
		}

		os.Exit(status)
	}

	return parser.OutArgs

}

func Parse(parameters interface{}, help string) []string {
	return HammerArguments(os.Args[0], os.Args[1:], parameters, help, true)
}
