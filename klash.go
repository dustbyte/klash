package klash

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
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

type ArgumentParser struct {
	Parser  *ParamParser
	Args    []string
	OutArgs []string
	Idx     int
	Stop    bool
	Stopped bool
}

func NewArgumentParser(parser *ParamParser, args []string, stop bool) *ArgumentParser {
	return &ArgumentParser{
		parser,
		args,
		make([]string, 0, len(args)),
		0,
		stop,
		false,
	}
}

func (ap *ArgumentParser) Terminated() bool {
	return ap.Idx >= len(ap.Args)
}

func (ap *ArgumentParser) extractVal(stringval string, value *reflect.Value) error {
	switch value.Kind() {
	case reflect.String:
		value.Set(reflect.ValueOf(stringval))
	case reflect.Int:
		val, err := strconv.ParseInt(stringval, 0, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int(val)))
	case reflect.Uint:
		val, err := strconv.ParseUint(stringval, 0, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint(val)))
	case reflect.Float32:
		val, err := strconv.ParseFloat(stringval, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(float32(val)))
	case reflect.Float64:
		val, err := strconv.ParseFloat(stringval, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(float64(val)))
	default:
		return fmt.Errorf("klash: Invalid type %s",
			value.Kind(),
		)
	}
	return nil
}

func (ap *ArgumentParser) ParseOne() error {
	arg := ap.Args[ap.Idx]

	if ap.Stopped || arg[0] != '-' {
		ap.OutArgs = append(ap.OutArgs, arg)
		if ap.Stop {
			ap.Stopped = true
		}
		ap.Idx++
		return nil
	}

	arg = strings.ToLower(arg)

	for len(arg) > 0 && arg[0] == '-' {
		arg = arg[1:]
	}

	if param, ok := ap.Parser.Params[arg]; ok {
		if param.Value.Kind() == reflect.Bool {
			param.Value.Set(reflect.ValueOf(true))
		} else {
			ap.Idx++
			stringval := ap.Args[ap.Idx]

			if param.Value.Kind() == reflect.Slice {
				value := reflect.New(param.Value.Type().Elem()).Elem()
				if err := ap.extractVal(stringval, &value); err != nil {
					return err
				}
				param.Value.Set(reflect.Append(param.Value, value))
			} else if err := ap.extractVal(stringval, &param.Value); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("klash: Invalid flag: %s", arg)
	}

	ap.Idx++
	return nil
}

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
