package klash

import (
	"fmt"
	"os"
	"reflect"
)

func ParseArguments(name string, arguments []string, parameters interface{}, stop bool) ([]string, error) {
	pvalue := reflect.ValueOf(parameters)
	if pvalue.Kind() != reflect.Ptr || pvalue.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s: Pointer to struct expected", name)
	}

	params := NewParams()
	if err := params.Parse(&pvalue); err != nil {
		return nil, err
	}

	parser := NewArgumentParser(name, params, arguments, stop)

	for !parser.Terminated() {
		if err := parser.ParseOne(); err != nil {
			return nil, err
		}
	}

	return parser.OutArgs, nil
}

func Parse(parameters interface{}) ([]string, error) {
	return ParseArguments(os.Args[0], os.Args[1:], parameters, true)
}
