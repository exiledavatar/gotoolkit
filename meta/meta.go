package meta

import (
	"reflect"
)

// ToIndirectReflectValue attempts to convert any value to a reflect.Value and indirect it. If value is nil,
// an invalid reflect.Value, or a nil pointer, the returned value will be invalid. It returns the indrected value
// and whether the original value was a pointer
func ToIndirectReflectValue(value any) (reflect.Value, reflect.Type, bool) {
	var rv reflect.Value
	var rt reflect.Type

	switch v := value.(type) {
	case nil:
		// do nothing
	case reflect.Type:
		rv = reflect.New(v).Elem()
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	var pointer bool
	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = rv.Elem()
	default:
		rt = rv.Type()
	}

	return rv, rt, pointer
}
