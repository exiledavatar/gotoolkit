package meta

import (
	"reflect"
)

type Data []any

func ToData(value any) Data {
	// var data Data
	// for _, v := range value {
	// 	data = append(data, v)
	// }

	// return data
	if value == nil {
		return nil
	}
	// fmt.Printf("Value is %#v -----------------------------------\n", value)

	return Data(ToSlice(value))
	// return Data{}
}

// ToSlice is intended to explicity convert a slice to a slice of type any
func ToSlice(a any) []any {
	v, err := ToValue(a)
	if err != nil {
		return nil
	}
	var s []any
	switch {
	case v.Kind() == reflect.Invalid:
		return nil
	case v.Kind() == reflect.Slice || v.Kind() == reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s = append(s, v.Index(i).Interface())
		}
	case v.Kind() == reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			s = append(s, iter.Value().Interface())
		}
	default:
		s = append(s, a)
	}
	return s

}
