package klash

import (
	"errors"
	"os"
	"reflect"
)

func ParseArguments(arguments []string, parameters interface{}, stop bool) ([]string, error) {
	pvalue := reflect.ValueOf(parameters)
	if pvalue.Kind() != reflect.Ptr || pvalue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("klash: Pointer to struct expected")
	}

	params := NewParams()
	if err := params.Parse(&pvalue); err != nil {
		return nil, err
	}
	parser := NewArgumentParser(params, arguments, stop)

	for !parser.Terminated() {
		if err := parser.ParseOne(); err != nil {
			return nil, err
		}
	}

	return parser.OutArgs, nil
}

func Parse(parameters interface{}) ([]string, error) {
	return ParseArguments(os.Args[1:], parameters, true)
}
