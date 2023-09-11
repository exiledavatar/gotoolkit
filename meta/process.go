package meta

import (
	"fmt"
	"log"
	"reflect"
)

type Processor interface {
	Process() error
}

func IsProcessor(a any) bool {
	return ImplementsInterface[Processor](a)
}

func AsProcessor[T any](a T) (Processor, error) {
	if IsProcessor(a) {
		var iv interface{} = a
		return (iv).(Processor), nil
	}
	return nil, fmt.Errorf("%T does not implement Processor", a)
}

func Process(a any) error {
	switch v, err := ToValue(a); {
	case err != nil:
		return err
	case IsProcessor(a):
		return a.(Processor).Process()
	case v.Kind() == reflect.Struct:
		return ProcessStruct(a)
	case v.Kind() == reflect.Slice:
		return ProcessSlice(a)
	case v.Kind() == reflect.Map:
		return ProcessMap(a)
	default:
		return nil
	}
}

// ProcessRecursively dive's into any member or element processing.
// It then attempts to call a's Process method, if applicable.
func ProcessRecursively(a any) error {
	switch v, err := ToValue(a); {
	case err != nil:
		return err
	case v.Kind() == reflect.Struct:
		if err := ProcessStruct(a); err != nil {
			return err
		}
	case v.Kind() == reflect.Slice:
		if err := ProcessSlice(a); err != nil {
			return err
		}
	case v.Kind() == reflect.Map:
		if err := ProcessMap(a); err != nil {
			return err
		}
	default:
	}
	if IsProcessor(a) {
		return a.(Processor).Process()
	}
	return nil
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
