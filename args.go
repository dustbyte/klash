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
	Params  *Params
	Args    []string
	OutArgs []string
	Idx     int
	Stop    bool
	Stopped bool
}

func NewArgumentParser(params *Params, args []string, stop bool) *ArgumentParser {
	return &ArgumentParser{
		params,
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
	value *reflect.Value) (error, bool) {

	var ptrType reflect.Type
	var ptr reflect.Value

	if value.Kind() == reflect.Ptr {
		ptrType = value.Type()
		ptr = *value
	} else {
		ptrType = reflect.PtrTo(value.Type())

		if !value.CanAddr() {
			return fmt.Errorf("%s is not addressable", value.Type()), false
		}
		ptr = value.Addr()
	}

	interfaceType := reflect.TypeOf((*Convertible)(nil)).Elem()

	if !ptrType.Implements(interfaceType) {
		return nil, false
	}

	method := ptr.MethodByName("FromString")
	if !method.IsValid() {
		return fmt.Errorf("conversion method not valid: %s", value.Type()), false
	}

	ierr := method.Call([]reflect.Value{reflect.ValueOf(stringval)})[0].Interface()

	// Cannot assert nil to be of type error
	if ierr == nil {
		return nil, true
	}
	return (ierr).(error), false
}

func (ap *ArgumentParser) extractVal(stringval string,
	value *reflect.Value) error {

	err, converted := ap.checkConvertible(stringval, value)

	if err != nil {
		return err
	}

	if !converted {
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
			return fmt.Errorf("cannot handle type %s", value.Type())
		}
	}
	return nil
}

func (ap *ArgumentParser) setBool(param *Parameter) {
	param.Value.Set(reflect.ValueOf(true))
}

func (ap *ArgumentParser) stripDashes(arg string) string {
	for len(arg) > 0 && arg[0] == '-' {
		arg = arg[1:]
	}
	return arg
}

func (ap *ArgumentParser) explodeArg(arg string) (string, string, error) {
	idx := strings.Index(arg, "=")

	if idx >= 0 {
		exploded := strings.Split(arg, "=")
		if exploded[1] == "" {
			return "", "", fmt.Errorf("no value provided to %s", exploded[0])
		}
		return exploded[0], exploded[1], nil
	}
	return arg, "", nil
}

func (ap *ArgumentParser) ParseOne() error {
	arg := ap.Args[ap.Idx]
	var stringval string
	var err error

	if ap.Stopped || arg[0] != '-' {
		ap.OutArgs = append(ap.OutArgs, arg)
		if ap.Stop {
			ap.Stopped = true
		}
		ap.Idx++
		return nil
	}

	arg = ap.stripDashes(arg)
	arg, stringval, err = ap.explodeArg(arg)
	if err != nil {
		return err
	}

	arg = strings.ToLower(arg)

	if arg == "help" || arg == "h" {
		return HelpError
	}

	if param, ok := ap.Params.Get(arg); ok {
		if param.Value.Kind() == reflect.Bool {
			ap.setBool(param)
		} else {
			if stringval == "" {
				ap.Idx++
				if ap.Idx >= len(ap.Args) {
					return fmt.Errorf("no value provided for %s", arg)
				}
				stringval = ap.Args[ap.Idx]
			}

			if param.Value.Kind() == reflect.Slice {
				value := reflect.New(param.Value.Type().Elem()).Elem()
				if err = ap.extractVal(stringval, &value); err != nil {
					return err
				}
				param.Value.Set(reflect.Append(param.Value, value))
			} else if err = ap.extractVal(stringval, &param.Value); err != nil {
				return err
			}
		}
	} else {
		for _, rune := range arg {
			param, ok := ap.Params.Get(string(rune))
			if ok && param.Value.Kind() == reflect.Bool {
				ap.setBool(param)
			} else {
				return fmt.Errorf("unrecognized arguments: %s", arg)
			}
		}
	}

	ap.Idx++
	return nil
}

func (ap *ArgumentParser) Parse() error {
	for !ap.Terminated() {
		if err := ap.ParseOne(); err != nil {
			return err
		}
	}
	return nil
}
