package klash

import (
	"reflect"
	"strings"
)

// A Parameter is a convenient representation of an expected parameter.
// Since parameters can have aliases (e.g -v and --version), a Parameter
// can store multiple (expected to be around 2) names.
// The Value field is the reflection representation of the
// value of the parameters structure field.
type Parameter struct {
	Name  string
	Alias string
	Value reflect.Value
	Help  string
}

// A capacity of 2 seems to be a good guess for the number of aliases.
func NewParameter(name string, value reflect.Value) *Parameter {
	parameter := Parameter{
		Name:  name,
		Value: value,
	}
	return &parameter
}

// List through all Tag-related (i.e starting with Tag) method of the Parameter type
// and execute them.
func (p *Parameter) DiscoverProperties(tag reflect.StructTag) error {
	if len(tag) > 0 {
		paramtype := reflect.TypeOf((*Parameter)(nil))
		prefix := "Tag"
		paramvalue := reflect.ValueOf(p)

		for idx := 0; idx < paramtype.NumMethod(); idx++ {
			method := paramtype.Method(idx)
			if !strings.HasPrefix(method.Name, prefix) {
				continue
			}

			tagname := "klash-" + strings.ToLower(method.Name[len(prefix):])
			if tagval := tag.Get(tagname); tagval != "" {
				methodValue := paramvalue.MethodByName(method.Name)

				err := methodValue.Call([]reflect.Value{reflect.ValueOf(tagval)})[0].Interface()
				if err != nil {
					return (err).(error)
				}
			}
		}
	}
	return nil
}

func (p *Parameter) TagHelp(tag string) error {
	p.Help = tag
	return nil
}

func (p *Parameter) TagAlias(tag string) error {
	p.Alias = tag
	return nil
}
