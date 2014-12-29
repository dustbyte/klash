package klash

import (
	"fmt"
	"reflect"
	"strings"
)

var allowedTypes map[reflect.Kind]bool = map[reflect.Kind]bool{
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Uint:    true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  true,
	reflect.Slice:   true,
}

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
		if _, ok := allowedTypes[value.Kind()]; !ok {
			return fmt.Errorf("klash: Invalid type for parameter %s: %s",
				field.Name,
				value.Kind(),
			)
		}

		if value.Kind() == reflect.Slice {
			sliceType := value.Type().Elem()
			_, ok := allowedTypes[sliceType.Kind()]
			if !ok || sliceType.Kind() == reflect.Slice {
				return fmt.Errorf("klash: Invalid slice type for parameter %s: %s",
					field.Name,
					sliceType.Kind(),
				)
			}
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
		}

		parameter := NewParameter(field.Name, value)

		for _, name := range parameter.Names {
			p.Params[strings.ToLower(name)] = parameter
		}
	}
	return nil
}
