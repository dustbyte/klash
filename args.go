package klash

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Convertible interface {
	FromString(stringval string) error
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

func (ap *ArgumentParser) checkConvertible(stringval string,
	value *reflect.Value) error {

	if !value.CanAddr() {
		return fmt.Errorf("klash: %s is not addressable", value.Type())
	}
	ptr := value.Addr()

	interfaceType := reflect.TypeOf((*Convertible)(nil)).Elem()

	if !ptr.Type().Implements(interfaceType) {
		return fmt.Errorf("klash: Invalid type %s", value.Type())
	}

	method := ptr.MethodByName("FromString")
	if !method.IsValid() {
		return fmt.Errorf("klash: Method not valid for %s", value.Type())
	}

	ierr := method.Call([]reflect.Value{reflect.ValueOf(stringval)})[0].Interface()

	// Cannot assert nil to be of type error
	if ierr == nil {
		return nil
	}
	return (ierr).(error)
}

func (ap *ArgumentParser) extractVal(stringval string,
	value *reflect.Value) error {

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
		return ap.checkConvertible(stringval, value)
	}
	return nil
}

func (ap *ArgumentParser) setBool(param *Parameter) {
	param.Value.Set(reflect.ValueOf(true))
}

func (ap *ArgumentParser) ParseOne() error {
	arg := ap.Args[ap.Idx]
	var stringval string

	if ap.Stopped || arg[0] != '-' {
		ap.OutArgs = append(ap.OutArgs, arg)
		if ap.Stop {
			ap.Stopped = true
		}
		ap.Idx++
		return nil
	}

	for len(arg) > 0 && arg[0] == '-' {
		arg = arg[1:]
	}

	idx := strings.Index(arg, "=")
	if idx >= 0 {
		exploded := strings.Split(arg, "=")
		if exploded[1] == "" {
			return fmt.Errorf("klash: no value provided to %s", exploded[0])
		}
		arg, stringval = exploded[0], exploded[1]
	}

	arg = strings.ToLower(arg)

	if param, ok := ap.Parser.Params[arg]; ok {
		if param.Value.Kind() == reflect.Bool {
			ap.setBool(param)
		} else {
			if stringval == "" {
				ap.Idx++
				if ap.Idx >= len(ap.Args) {
					return fmt.Errorf("klash: no value provided to %s", arg)
				}
				stringval = ap.Args[ap.Idx]
			}

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
		for _, rune := range arg {
			param, ok := ap.Parser.Params[string(rune)]
			if ok && param.Value.Kind() == reflect.Bool {
				ap.setBool(param)
			} else {
				return fmt.Errorf("klash: Invalid flag: %s", arg)
			}
		}
	}

	ap.Idx++
	return nil
}
