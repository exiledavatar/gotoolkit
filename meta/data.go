package meta

import (
	"fmt"
	"reflect"
)

type Data []any

func ToData(value any) Data {
	// var data Data
	// for _, v := range value {
	// 	data = append(data, v)
	// }

	// return data
	return Data(ToSlice(value))
}

// ToSlice is intended to explicity convert a slice to a slice of type any
func ToSlice(a any) []any {

	v, err := ToValue(a)
	if err != nil {
		panic(err)
	}
	var s []any
	switch {
	case v.Kind() == reflect.Invalid:
		panic(fmt.Sprintf("%v is invalid", a))
	case v.Kind() == reflect.Slice || v.Kind() == reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s = append(s, v.Index(i).Interface())
		}
	case v.Kind() == reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			s = append(s, iter.Value().Interface())
		}
	case v.Kind() == reflect.Struct:
		s = append(s, v.Interface())
	}
	return s

}
