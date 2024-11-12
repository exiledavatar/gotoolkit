package meta

import (
	"fmt"
	"reflect"
	"slices"
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

type StructWithData struct {
	Struct
	Data
}

// ToAnySlice converts a slice []T to []any by wrapping each element in interfaces. Useful when you have a concrete
// slice and need to feed it to something that requires []any.
func ToAnySlice[T any](value []T) []any {
	var anyValue []any
	for _, v := range value {
		anyValue = append(anyValue, v)
	}
	return anyValue
}

func Hack(values ...any) []any {
	var out []any

	return out
}

// ToSlice should convert a mix of single values and slices to []any, just don't
// get too tricky
func ToSlice(values ...any) []any {
	var out []any
	for _, value := range values {
		switch rv := reflect.ValueOf(value); {
		case !rv.IsValid():
		// do nothing
		case rv.Kind() == reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				if rvi := rv.Index(i); rvi.IsValid() {
					out = append(out, rvi.Interface())
				}

			}
		default:
			out = append(out, rv.Interface())
		}
	}
	return out

}

func RemoveNils(values []any) []any {
	var out []any
	for _, v := range values {
		if v != nil {
			out = append(out, v)
		}
	}
	return out
}

// Flatten is intended for ragged slices that potentially include slices of slices (of slices of slices...)
// Don't push your luck
func Flatten(values any) []any {
	var out []any
	switch rv := reflect.ValueOf(values); {
	case !rv.IsValid():
		return nil
	case rv.Kind() == reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if rvi := rv.Index(i); rvi.IsValid() {
				rvif := Flatten(rvi.Interface())
				if rvif != nil {

					out = slices.Concat(out, rvif)
				}
			}
		}
	default:
		out = append(out, rv.Interface())
	}
	return out
}

func ToStringSlice(values ...any) []string {
	out := []string{}
	vs := ToSlice(values...)
	vs = Flatten(vs)
	vs = Flatten(vs)
	for _, v := range vs {
		out = append(out, fmt.Sprint(v))
	}
	return out
}

// Coalesce tries to mimic similar sql functions: it iterates through slices at the same time and
// returns a slice with the first non-zero value for each element. Unequal length slices are not an issue
func Coalesce[S []T, T any](first, second S, dflt T) S {
	out := S{}
	lf := len(first)
	ls := len(second)
	lmax := max(lf, ls)

	for i := 0; i < lmax; i++ {
		var fi, si T
		if lf > i {
			fi = first[i]
		}
		if ls > i {
			si = second[i]
		}
		fiv := reflect.ValueOf(fi)
		siv := reflect.ValueOf(si)
		switch {
		case fiv.IsValid() && !fiv.IsZero():
			out = append(out, fi)
		case siv.IsValid() && !siv.IsZero():
			out = append(out, si)
		default:
			out = append(out, dflt)
		}
	}
	return out
}
