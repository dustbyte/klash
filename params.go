package klash

import "reflect"

// A Parameter is a convenient representation of an expected parameter.
// Since parameters can have aliases (e.g -v and --version), a Parameter
// can store multiple (expected to be around 2) names.
// The Value field is the reflection representation of the
// value of the parameters structure field.
type Parameter struct {
	Names []string
	Value reflect.Value
}

// The capacity of 2 seems to be a good guess for the number of aliases.
func NewParameter(name string, value reflect.Value) *Parameter {
	parameter := Parameter{
		make([]string, 1, 2),
		value,
	}
	parameter.Names[0] = name
	return &parameter
}

type Params map[string]*Parameter

// Params store the mapping of ParamName -> Parameter for the given structure.
// Since multiple names can be affected to a single parameter, multiple
// keys can be associated with a single parameter.
func NewParams() Params {
	return make(map[string]*Parameter)
}

// Parse discovers the given parameters structure and associates the structure's
// field names with their values into the Params structure.
func (p Params) Parse(pvalue *reflect.Value) {
	vtype := pvalue.Type().Elem()

	for idx := 0; idx < vtype.NumField(); idx++ {
		field := vtype.Field(idx)

		value := pvalue.Elem().Field(idx)

		if value.Kind() == reflect.Slice {
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
		}

		parameter := NewParameter(field.Name, value)

		for _, name := range parameter.Names {
			p[DecomposeName(name)] = parameter
		}
	}
}
