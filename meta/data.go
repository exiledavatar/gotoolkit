package meta

import (
	"log"
	"reflect"
)

type Data []any

func ToData[T any](value []T) Data {
	var data Data
	for _, v := range value {
		data = append(data, v)
	}

	return data
}

// ToSlice is intended to explicity convert a slice to a slice of type any
func ToSlice(a any) []any {
	v := reflect.ValueOf(a)
	var s []any
	if v.Kind() != reflect.Slice {
		log.Printf("%T is not a slice\n", a)
	}
	if v.Kind() == reflect.Slice {
		iv := reflect.Indirect(v)
		sliceType := reflect.TypeOf(a).Elem()
		out := reflect.MakeSlice(reflect.SliceOf(sliceType), iv.Len(), iv.Len())
		for i := 0; i < iv.Len(); i++ {
			f := reflect.Indirect(iv.Index(i))
			out.Index(i).Set(f)
			s = append(s, f.Interface())
		}
	}
	return s
}
