package meta

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// Value is an attempt to capture a dereferenced reflect.Value.
// let's hope it works
type Value struct {
	Name string // optional - intended for struct field names, map keys, or slice/array indexes
	reflect.Value
	Pointer bool
	Parent  *Value // only populated if it is a member of a struct or elements of collections - slice, array, map, or channel
}

// ToValue should be able to return a valid Value for anything except a nil pointer or an invalid reflect.Value
func ToValue(value any) (Value, error) {
	var v Value
	rv, rt, pointer := ToIndirectReflectValue(value)

	switch {
	case rt == nil:
		return v, fmt.Errorf("invalid value: %v", value)
	case rt == reflect.TypeOf(Value{}):
		return value.(Value), nil
	default:
		v = Value{
			Name:    rt.String(), // default name is the package.name of the type
			Value:   rv,
			Pointer: pointer,
		}
		return v, nil
	}

}

func (v Value) Kind() reflect.Kind {
	return v.Value.Kind()
}

// Valid just wraps reflect.Value.IsValid, because I don't like IsXXX methods unless I really really have to
func (v Value) Valid() bool {
	return v.Value.IsValid()
}

func (v Value) TypeMap(system string) string {
	return TypeMappings.To(system, v.Interface())
}

// Children returns the elements of slices, arrays, and maps, and the non-anonymous, exported fields of structs
// if it can be considered a 'child' in some way, it should be returned by this method
func (v *Value) Children() []Value {
	var children []Value
	switch kind := v.Kind(); {
	case kind == reflect.Invalid:
		return children
	case kind == reflect.Chan:
		// treat channel's as 0 length
		fallthrough
	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind) &&
		v.Value.Len() == 0:

		return children
	case kind == reflect.Map:
		iter := v.Value.MapRange()
		for iter.Next() {
			child, err := ToValue(iter.Value().Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return children
			}
			child.Name = iter.Key().String()
			// children[iter.Key().String()] = child
			children = append(children, child)
		}
	case kind == reflect.Slice || kind == reflect.Array:
		// precision := len(fmt.Sprint(rv.Len())) + 2
		for i := 0; i < v.Value.Len(); i++ {
			child, err := ToValue(v.Value.Index(i).Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return children
			}
			// key := fmt.Sprintf("%.[1]*d\n", precision, i)
			// children[key] = child
			// child.Name = key
			children = append(children, child)
		}
	case kind == reflect.Struct:
		for _, field := range reflect.VisibleFields(v.Type()) {
			if field.Anonymous || !field.IsExported() {
				continue
			}
			fieldValue := v.Value.FieldByName(field.Name)
			if fieldValue.Kind() == reflect.Invalid {
				fieldValue = reflect.New(field.Type).Elem()
			}
			// fmt.Println(field.Name, ":\t", fieldValue.Type())
			child, err := ToValue(fieldValue.Interface())
			if err != nil {
				return children
			}
			// children[field.Name] = child
			child.Name = field.Name
			children = append(children, child)
		}
	default:
	}

	for _, child := range children {
		child.Parent = v
	}
	return children

}

// Child will attempt to return a value's child with the given index/key/field name.
// It accepts an int for slices, arrays, and structs.
// It accepts a string for maps and structs.
func (v *Value) Child(a any) (Value, error) {
	children := v.Children()
	if len(children) == 0 {
		return Value{}, fmt.Errorf("%v has no children", v)
	}
	switch t := a.(type) {
	case int:
		if slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Struct}, v.Kind()) && len(children) >= t {
			return children[t], nil
		}
		// return Value{}, nil
	case string:
		if slices.Contains([]reflect.Kind{reflect.Map, reflect.Struct}, v.Kind()) {
			for _, child := range children {
				if child.Name == t {
					return child, nil
				}
			}
		}
	default:
		return Value{}, fmt.Errorf("method child(a any) expects an int or string, got %v", t)
	}
	return Value{}, fmt.Errorf("no matching child for %v", a)
}

func (v *Value) ChildrenByIndex(index ...any) ([]Value, error) {
	var children []Value
	for _, ind := range index {
		child, err := v.Child(ind)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

// NewElement returns a blank element for slices, arrays, maps and channels.
// For everything else it returns an (invalid) Value and a non-nil error.
func (v Value) NewElement() (Value, error) {
	var element Value
	if !v.Valid() {
		return element, fmt.Errorf("cannot determine child type of invalid value")
	}
	if slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, v.Kind()) {
		elemRT := v.Value.Type().Elem()
		return ToValue(elemRT)
	}
	return element, fmt.Errorf("value %v is a %s, not slice, array, map, or channel", v, v.Kind())
}
