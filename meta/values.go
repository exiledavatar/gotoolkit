package meta

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type Values []Value

func ToValues(a ...any) (Values, error) {
	var values []Value
	for _, ai := range a {
		value, err := ToValue(ai)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func (v Values) Names() []string {
	var names []string
	for _, value := range v {
		names = append(names, value.Name)
	}
	return names
}

func (v Values) ByNames(names ...string) Values {
	var values []Value
	for i, valueName := range v.Names() {
		if slices.Contains(names, valueName) {
			values = append(values, v[i])
		}
	}
	return values
}

func (v Values) Types() []reflect.Type {
	var types []reflect.Type
	for _, value := range v {
		types = append(types, value.Type())
	}
	return types
}

func (v Values) ByTypes(types ...reflect.Type) Values {
	var values []Value
	for i, valueType := range v.Types() {
		if slices.Contains(types, valueType) {
			values = append(values, v[i])
		}
	}
	return values
}
