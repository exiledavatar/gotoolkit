package processor

import (
	"log"
	"reflect"
)

type Processor interface {
	Process() error
}

// IsType is a convenience for use in switch statements
func IsType[A any](value any) bool {
	_, ok := value.(A)
	return ok
}

func IsString(value any) bool {
	return IsType[string](value)
}

// IsProcessor is for convenience, primarily for switch statements
func IsProcessor(a any) bool {
	_, ok := a.(Processor)
	return ok
}

// // AsProcessor uses assertion to attempt to return any value as a processor
// func AsProcessor[T any](a T) (Processor, error) {
// 	if IsProcessor(a) {
// 		var iv interface{} = a
// 		return (iv).(Processor), nil
// 	}
// 	return nil, fmt.Errorf("%T does not implement Processor", a)
// }

// IsKind returns true if a's reflect.Kind == t
func IsKind(a any, kind reflect.Kind) bool {
	return IndirectReflectionValue(a).Kind() == kind
}

func Process(a any) error {
	p, ok := a.(Processor)
	if ok {
		return p.Process()
	}
	switch {
	case IsProcessor(a):
		err := a.(Processor).Process()
		// fmt.Printf("Processed as a processor - results: %#v\n", a)
		return err
	case IsKind(a, reflect.Struct):
		return ProcessStruct(a)
	case IsKind(a, reflect.Slice):
		return ProcessSlice(a)
	case IsKind(a, reflect.Map):
		return ProcessMap(a)
	default:
		// fmt.Printf("Process(%s %T): no defined processing path - doing nothing.\n", reflect.ValueOf(a).Kind(), a)
		// Process ignores any fields that don't need processing
		return nil
	}
	// return errors.New(fmt.Sprintf("%T is not a Processor", a))
}

// ProcessRecursively dive's into any member or element processing.
// It then attempts to call a' Process method, if applicable.
func ProcessRecursively(a any) error {
	// v := IndirectReflectionValue(a)
	// fmt.Printf("%#v\n", v)
	switch {
	case IsKind(a, reflect.Struct):
		if err := ProcessStruct(a); err != nil {
			return err
		}
	case IsKind(a, reflect.Slice):
		if err := ProcessSlice(a); err != nil {
			return err
		}
	case IsKind(a, reflect.Map):
		if err := ProcessMap(a); err != nil {
			return err
		}
	default:
		// Process ignores any fields that dont' need processing
		// return nil
	}
	if IsProcessor(a) {
		if err := a.(Processor).Process(); err != nil {
			return err
		}
	}
	return nil
	// return errors.New(fmt.Sprintf("%T is not a Processor", a))
}

// ProcessSlice will only attempt to process elements, it won't
// manage the slice header itself - ie - it will modify or set elements
// to nil, but will not attempt to grow
func ProcessSlice(a any) error {
	v := reflect.ValueOf(a)
	iv := reflect.Indirect(v)
	switch iv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < iv.Len(); i++ {
			f := iv.Index(i)
			fv := f.Interface()
			if err := Process(fv); err != nil {
				// return err
				log.Println(err)
			}
		}
	default:
		log.Printf("%s is not a slice", v.Kind())
	}
	// fmt.Println("I'm a slice")
	return nil
}

// ProcessMap is used for dispatching - prefer Process as it will call
// ProcessMap for any map type that doesn't implement Processor interface
func ProcessMap(a any) error {
	v := reflect.ValueOf(a)
	iv := reflect.Indirect(v)
	if iv.Kind() != reflect.Map {
		log.Printf("%s is not a map", iv.Kind())
		return nil
	}
	iter := iv.MapRange()
	for iter.Next() {
		miv := iter.Value()
		mv := reflect.Indirect(miv)
		temp := reflect.New(mv.Type())
		temp.Elem().Set(mv)
		if err := Process(temp.Interface()); err != nil {
			// return err
			log.Println(err)
			return nil
		}
		if miv.Kind() == reflect.Pointer {
			iv.SetMapIndex(iter.Key(), temp)
		} else {
			iv.SetMapIndex(iter.Key(), temp.Elem())

		}
	}
	return nil
}

// ProcessStruct doesn't attempt to check/use the struct's Process method.
//
//	Instead it iterates through each member and attempts to Process them.
//
// It also makes no effort to process members that are nil pointers or
// otherwise result in reflect.Kind() == reflect.Invalid.
func ProcessStruct(a any) error {
	v := reflect.ValueOf(a)
	iv := reflect.Indirect(v)
	// fields := StructFields(iv)
	var fields []reflect.StructField
	for i := 0; i < iv.Type().NumField(); i++ {
		fields = append(fields, iv.Type().Field(i))
	}

	// fmt.Printf("ProcessStruct %s %s\n", iv.Type(), iv.Kind())

	for i := 0; i < iv.NumField(); i++ {
		sf := fields[i]
		// fmt.Printf("ProcessStruct Field %d %s\n", i, sf.Name)
		fv := iv.Field(i)
		fiv := reflect.Indirect(fv)
		switch {
		case !sf.IsExported() || fv.IsZero():
			// skip unexported fields or invalid fields
			continue
		case v.Kind() == reflect.Pointer && fv.Kind() == reflect.Pointer:
			// fmt.Printf("%s %s : %s %s : %s\n", v.Kind(), v.Type(), fv.Kind(), fv.Type(), sf.Name)
			if err := Process(fv.Interface()); err != nil {
				log.Println(err)
			}
		case v.Kind() == reflect.Pointer && fv.Kind() != reflect.Pointer:
			ptrValue := reflect.New(fiv.Type())
			ptrValue.Elem().Set(fiv)
			// fmt.Printf("%s %s : %s %s : %s\n", v.Kind(), v.Type(), fv.Kind(), fv.Type(), sf.Name)
			if err := Process(ptrValue.Interface()); err != nil {
				log.Println(err)
			}
			fv.Set(ptrValue.Elem())
		case v.Kind() != reflect.Pointer && fv.Kind() == reflect.Pointer:
			// this might not work - might need to follow pattern above
			if err := Process(fv.Interface()); err != nil {
				log.Println(err)
			}

		case v.Kind() != reflect.Pointer && fv.Kind() != reflect.Pointer:
			// do nothing - can't use processor interface, like, at all...
			continue
		default:
			continue
		}
	}
	return nil
}

// IndirectReflectionValue attempts to convert a to
// an indirect reflection value, dereference it, and return it
func IndirectReflectionValue(a any) reflect.Value {
	rv, ok := a.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(a)
	}
	if rv.Kind() == reflect.Pointer {
		return reflect.Indirect(rv)
	}
	return rv
}

// func ProcessStructFields(a any) {
// 	fmt.Println("--------------------------------------------------------------------------")
// 	v := reflect.ValueOf(a)
// 	fmt.Printf("%T\n", v)
// 	vi := reflect.Indirect(v)
// 	fmt.Printf("%T\n", vi)
// 	// reflect.ValueOf(&t).MethodByName("GFG").Call([]reflect.Value{})
// 	fields := StructFields(a)
// 	// for _, f := range fields {
// 	for i := 0; i < v.NumField(); i++ {
// 		// fmt.Println("--------------------------------------")
// 		// sf := v.Type().Field(i)
// 		f := v.Field(i) //.Addr()
// 		f = reflect.Indirect(f)
// 		if f.Kind() == reflect.Invalid || !f.CanSet() {
// 			fmt.Printf("%#v\n", fields[i])
// 		}
// 		if f.Kind() == reflect.Slice {
// 			// f = f.Elem()
// 			// fmt.Printf("%#v\n", f)
// 			ProcessSlice(f)
// 		}
// 		fmt.Printf("%T:  %s is settable: %t\n", a, f.Kind(), f.CanSet())
// 		fmt.Println("--------------------------------------")
// 		// fmt.Printf("%#v\n", f)
// 		// if f.Kind() == reflect.Pointer {
// 		// 	fv := f.Elem()
// 		// 	fv = reflect.Indirect(fv)
// 		// 	// fmt.Printf("%s: %#v\n", sf.Name, fv)

// 		// 	// if f.IsNil() {
// 		// 	// 	fmt.Printf("%T.%s is nil but settable: %t\n", a, f.Type(), fv.CanSet())
// 		// 	// 	// fmt.Println(f.Addr())
// 		// 	// 	fmt.Println("--------------------------------------")
// 		// 	// 	continue
// 		// 	// }
// 		// 	fmt.Printf("%T.%s is settable: %t\n", a, fv.Type(), fv.CanSet())
// 		// }
// 		// fmt.Printf("%T.%s is settable: %t\n", a, f.Type(), (reflect.Indirect(f)).CanSet())
// 		// fmt.Println("--------------------------------------")

// 		// fmt.Printf("%T is pointer: %t\n", a, f.Kind() == reflect.Pointer)

// 		// f.MethodByName("Process").Call([]reflect.Value{})
// 		// fv := v.FieldByName(f.Name)

// 		// fmt.Printf("%#v\n", fv)
// 		// res := fv.MethodByName("Process").Call([]reflect.Value{})
// 		// fmt.Println(res)
// 		// fmt.Printf("%#v\n", f)

// 	}
// 	fmt.Printf("%s\n", vi.Type())
// }
