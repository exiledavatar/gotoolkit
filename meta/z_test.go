package meta_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/ericsgagnon/qgenda/pkg/meta"
)

// func TestZZZToMeta(t *testing.T) {
// 	fmt.Printf("%-.200s\n", strings.Repeat("-", 200))
// 	fmt.Println("ToMeta")
// 	t.Run("map[string]struct", func(t *testing.T) {
// 		v := map[string]ExampleStruct{}
// 		m := meta.ToMeta(v)
// 		fmt.Printf("\n%+v\n", m)
// 	})
// 	fmt.Printf("%-.200s\n", strings.Repeat("-", 200))
// 	x1 := "testValue"
// 	x := map[any]any{
// 		"testKey":     x1,
// 		"testPointer": &x1,
// 		"nilKey":      nil,
// 	}
// 	rv := reflect.ValueOf(x)
// 	fmt.Println(rv.Type())
// 	if rv.Kind() == reflect.Map {
// 		rv0 := rv.MapIndex(rv.MapKeys()[0])
// 		fmt.Println("----------")
// 		fmt.Println(rv0.Interface())
// 		meta.Unbox(rv0.Interface())
// 		fmt.Println("----------")
// 		if rv0.Kind() == reflect.Interface {
// 			rv0 = rv0.Elem()
// 			fmt.Println(rv0.Type())
// 			fmt.Println(rv0.Kind())

// 		}
// 		rv1 := rv.MapIndex(rv.MapKeys()[1])
// 		if rv1.Kind() == reflect.Interface {
// 			rv1 = rv1.Elem()
// 			fmt.Println(rv1.Kind())
// 		}
// 		meta.Unbox(x["testKey"])

// 		fmt.Println("meta.ToIndirect ----------")
// 		var boxed any
// 		boxed = x1
// 		rv, rt, pointer, bv := meta.ToIndirect(boxed)
// 		fmt.Println(rv, rt, pointer, bv)
// 		fmt.Println(meta.ToIndirect(any(boxed)))
// 		fmt.Println(meta.ToIndirect(&boxed))
// 		fmt.Println(reflect.ValueOf(boxed))
// 		fmt.Println("--------------------------")

// 	}
// 	fmt.Printf("%-.200s\n", strings.Repeat("-", 200))

// }

func TestToValue(t *testing.T) {
	t.Run("ToValue", func(t *testing.T) {
		v, err := meta.ToValue(structExample)
		if err != nil {
			log.Println(err)

		}
		// w := new(tabwriter.Writer)
		// w.Init(os.Stdout, 0, 8, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
		// // maxTab := 50
		// // toString := fmt.Sprintf("{}.{}")
		// precision := 30
		// columns := []any{"Type", "Value", "Pointer", "Parent", "Children"}
		// strings.Join()
		// header := fmt.Sprintf("", precision, columns...)
		// fmt.Fprintln(w, "Type\tValue\tPointer\tParent\tChildren")

		// fmt.Fprintln(w, "aaaaaaaaaaaaaaaaaaaa\tbbbbbbbbbbbbbbbbbbb\tccccccccccccccccccccccc\tdddddddddddddddddd\teeeeeeeeeee")
		// fmt.Fprintln(w)
		// w.Flush()
		for k, child := range v.Children() {
			fmt.Printf("%v\t%v\n", k, child)

		}
		fmt.Println("\n------------------------------------------------------------------------------")
		vx, err := meta.ToValue(structExample)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(vx)
		fmt.Println(vx.Children())
		fmt.Println("\n------------------------------------------------------------------------------")
		ch, err := vx.Child("ExampleStructSlice")
		if err != nil {
			log.Println(err)
		}
		fmt.Println("\n------------------------------------------------------------------------------")
		chElem, err := ch.NewElement()
		if err != nil {
			log.Println(err)
		}

		fmt.Println(chElem)
	})
}
