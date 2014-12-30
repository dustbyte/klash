package klash

import (
	"reflect"
	"strings"
)

type Parameter struct {
	Names []string
	Value reflect.Value
}

func NewParameter(name string, value reflect.Value) *Parameter {
	parameter := Parameter{
		make([]string, 1, 2),
		value,
	}
	parameter.Names[0] = name
	return &parameter
}

type ParamParser struct {
	Params map[string]*Parameter
}

func NewParamParser() *ParamParser {
	return &ParamParser{make(map[string]*Parameter)}
}

func (p *ParamParser) Parse(pvalue *reflect.Value) error {
	vtype := pvalue.Type().Elem()

	for idx := 0; idx < vtype.NumField(); idx++ {
		field := vtype.Field(idx)

		value := pvalue.Elem().Field(idx)

		if value.Kind() == reflect.Slice {
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
		}

		parameter := NewParameter(field.Name, value)

		for _, name := range parameter.Names {
			p.Params[strings.ToLower(name)] = parameter
		}
	}
	return nil
}
