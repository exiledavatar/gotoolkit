package meta

import (
	"reflect"
)

// types of fields:
// single values
// collections of single values
// single struct
// collections of structs

// type Structs []Struct

// Structs represents a slice, array, map, or channel of Structs
type Structs struct {
	Name        string
	Type        reflect.Type
	Value       reflect.Value
	Pointer     bool
	IndexType   reflect.Type
	ElementType reflect.Type
	Index       any
	Structs     []Struct
}



// func ToStructs(value any) (Structs, error) {
// 	var s Structs
// 	structs := []Struct{}
// 	rv, rt, pointer := ToIndirectReflectValue(value)
// 	var iType, eType reflect.Type
// 	var index any

// 	switch kind := rt.Kind(); {
// 	case kind == reflect.Invalid:
// 		return s, fmt.Errorf("invalid value: Kind() == reflect.Invalid: %v", value)
// 	case kind == reflect.Chan:
// 		// treat channel's a 0 length
// 		fallthrough
// 	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind) &&
// 		rv.Len() == 0:
// 	case kind == reflect.Map:
// 		iter := rv.MapRange()
// 		ind := []string{}
// 		for iter.Next() {
// iter.Key().Type()
// 			str, err := ToStruct(iter.Value().Interface())
// 			if str.Value.Kind() != reflect.Invalid && err != nil {
// 				return s, err
// 			}
// 			// str.Name = iter.Key().String()
// 			ind = append(ind, iter.Key().String())
// 			structs = append(structs, str)
// 		}
// 		index = ind
// 	case kind == reflect.Slice || kind == reflect.Array:
// 		ind := []int{}
// 		for i := 0; i < rv.Len(); i++ {
// 			str, err := ToStruct(rv.Index(i).Interface())

// 			if str.Value.Kind() != reflect.Invalid && err != nil {
// 				return s, err
// 			}
// 			str.Name = i
// 			children = append(children, child)
// 		}
// 	case kind == reflect.Struct:
// 		for _, field := range reflect.VisibleFields(rt) {
// 			if field.Anonymous || !field.IsExported() {
// 				continue
// 			}
// 			fieldValue := rv.FieldByName(field.Name)
// 			if fieldValue.Kind() == reflect.Invalid {
// 				fieldValue = reflect.New(field.Type).Elem()
// 			}
// 			fmt.Println(field.Name, ":\t", fieldValue.Type())
// 			child, err := ToValue(fieldValue.Interface())
// 			if err != nil {
// 				return v, err
// 			}
// 			// children[field.Name] = child
// 			child.Name = field.Name
// 			children = append(children, child)
// 		}
// 	default:
// 	}
// 	v = Value{
// 		Type:     rt,
// 		Value:    rv,
// 		Pointer:  pointer,
// 		Children: children,
// 	}

// 	return s, nil
// }

// func toStructs(value any) (Structs, error) {
// 	var s Structs
// 	return s, nil
// }

// func ToStructs[S ~[]T, T any](value T) (Structs, error) {
// 	var s Structs
// 	rv, rt, pointer := ToIndirectReflectValue(value)
// 	if !slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, rv.Kind()) {
// 		return s, fmt.Errorf("value %T is a %v, need a slice, array, map, or channel", value, rv.Kind())
// 	}
// 	if rv.Len() < 1 {
// 		value = *new(T)
// 	}
// 	var ss []Struct

// 	switch kind := rv.Kind(); {
// 	case kind == reflect.Slice:
// 	case kind == reflect.Map:
// 		iter := rv.MapRange()
// 		for iter.Next() {
// 			// k := iter.Key()oi
// 			v := iter.Value()
// 			ssi, err := ToStruct(v.Interface())
// 			if err != nil {
// 				return s, err
// 			}
// 			ss = append(ss, ssi)
// 		}
// 	case kind == reflect.Chan:
// 	case kind == reflect.Array:
// 	default:
// 	}

// 	// if len(value) < 1 {
// 	// 	// fmt.Printf("%T\n", *new(T))
// 	// 	s0, err := ToStruct(*new(T))
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	return append(s, s0), nil
// 	// }
// 	for _, v := range value {
// 		// fmt.Printf("%T\t%s\t%s\n", value, value, reflect.ValueOf(value).IsValid())
// 		// fmt.Printf("%T\t%s\t%s\n", v, v, reflect.ValueOf(v).IsValid())

// 		si, err := ToStruct(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = append(s, si)
// 	}
// 	return s, nil
// }

// func NewStructs[S ~[]T, T any](value S, cfg StructConfig) (Structs, error) {
// 	s := []Struct{}
// 	if len(value) < 1 {
// 		s0, err := NewStruct(*new(T), cfg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return append(s, s0), nil
// 	}
// 	for _, v := range value {
// 		si, err := NewStruct(v, cfg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = append(s, si)
// 	}
// 	return s, nil
// }

// func (s Structs) Names() []string {
// 	var names []string
// 	for _, v := range s.Structs {
// 		names = append(names, v.Name)
// 	}
// 	return names
// }

// set children:
// name
// *parent
// uuid
// Handle child slices of structs
// handle tags that indicate child/not child
//
