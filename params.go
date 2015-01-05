package klash

import "reflect"

type Params map[string]*Parameter

// Params store the mapping of ParamName -> Parameter for the given structure.
// Since multiple names can be affected to a single parameter, multiple
// keys can be associated with a single parameter.
func NewParams() Params {
	return make(map[string]*Parameter)
}

// Parse discovers the given parameters structure and associates the structure's
// field names with their values into the Params structure.
func (p Params) Parse(pvalue *reflect.Value) error {
	vtype := pvalue.Type().Elem()

	for idx := 0; idx < vtype.NumField(); idx++ {
		field := vtype.Field(idx)

		value := pvalue.Elem().Field(idx)

		if value.Kind() == reflect.Slice {
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
		}

		parameter := NewParameter(field.Name, value)
		if err := parameter.DiscoverProperties(field.Tag); err != nil {
			return err
		}

		for _, name := range parameter.Names {
			p[DecomposeName(name)] = parameter
		}
	}
	return nil
}
