package klash

import (
	"errors"
	"os"
	"reflect"
)

func ParseArguments(arguments []string, params interface{}, stop bool) ([]string, error) {
	pvalue := reflect.ValueOf(params)
	if pvalue.Kind() != reflect.Ptr {
		return nil, errors.New("klash: Pointer type expected")
	}

	parser := NewParamParser()
	if err := parser.Parse(&pvalue); err != nil {
		return nil, err
	}
	aparser := NewArgumentParser(parser, arguments, stop)

	for !aparser.Terminated() {
		if err := aparser.ParseOne(); err != nil {
			return nil, err
		}
	}

	return aparser.OutArgs, nil
}

func Parse(params interface{}) ([]string, error) {
	return ParseArguments(os.Args[1:], params, true)
}
