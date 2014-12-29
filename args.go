package klash

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
