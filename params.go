package klash

import (
	"errors"
	"fmt"
	"reflect"
)

type Params struct {
	Mapping map[string]*Parameter
	Listing []*Parameter
}

// Params store the mapping of ParamName -> Parameter for the given structure.
// Since multiple names can be affected to a single parameter, multiple
// keys can be associated with a single parameter.
func MakeParams(fieldCount int) *Params {
	return &Params{
		make(map[string]*Parameter),
		make([]*Parameter, 0, fieldCount),
	}
}

func NewParams(parameters interface{}) (*Params, error) {
	pvalue := reflect.ValueOf(parameters)
	if pvalue.Kind() != reflect.Ptr || pvalue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("klash: Pointer to struct expected")
	}
	fieldCount := pvalue.Type().Elem().NumField()

	params := MakeParams(fieldCount)
	if err := params.Parse(&pvalue); err != nil {
		return nil, err
	}

	return params, nil
}

// Parse discovers the given parameters structure and associates the structure's
// field names with their values into the Params structure.
func (p *Params) Parse(pvalue *reflect.Value) error {
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

		if err := p.Set(parameter.Name, parameter); err != nil {
			return err
		}

		if parameter.Alias != "" {
			if err := p.Set(parameter.Alias, parameter); err != nil {
				return err
			}
		}
		p.Listing = append(p.Listing, parameter)
	}
	return nil
}

func (p *Params) Get(key string) (*Parameter, bool) {
	val, ok := p.Mapping[DecomposeName(key, true)]
	return val, ok
}

func (p *Params) Set(key string, value *Parameter) error {
	key = DecomposeName(key, true)
	_, ok := p.Mapping[key]
	if ok {
		return fmt.Errorf("klash: %s is already an argument or an alias", key)
	}
	p.Mapping[key] = value
	return nil
}
