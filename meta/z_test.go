package meta_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/exiledavatar/gotoolkit/meta"
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

func TestX(t *testing.T) {
	t.Run("X", func(t *testing.T) {
		s, err := meta.ToStruct(structExample)
		if err != nil {
			log.Println(err)

		}

		for _, child := range s.Fields() {
			fmt.Printf("%s\t%s\t%v\t\t\n", child.Name, child.Type(), child.Tags())
		}

		fmt.Println("TypeMaps:---------------------------------------------------")
		for _, child := range s.Fields() {

			fmt.Printf("%s\t%s\t%s\t%s\n", child.Name, child.TypeMap("postgres"), child.Type(), child.Kind())
		}

		for _, field := range s.Fields() {
			fmt.Println(field.Name, field.TagName("json", "pg", "db"))
		}

		// fmt.Println(meta.TypeMappings)
		// v, err := meta.ToValue(structExample)
		// if err != nil {
		// 	log.Println(err)
		// }
		// fmt.Println(v.)
		// for _, field := range ss.Fields() {
		// 	fmt.Println(field.Name, field.Type(), field.MultiValued())
		// }

		// for k, child := range s.Children() {
		// 	fmt.Printf("%v\t%v\n", k, child)

		// }
		// fmt.Println("\n------------------------------------------------------------------------------")
		// vx, err := meta.ToValue(structExample)
		// if err != nil {
		// 	log.Println(err)
		// }
		// fmt.Println(vx)
		// fmt.Println(vx.Children())
		// fmt.Println("\n------------------------------------------------------------------------------")
		// ch, err := vx.Child("ExampleStructSlice")
		// if err != nil {
		// 	log.Println(err)
		// }
		// fmt.Println("\n------------------------------------------------------------------------------")
		// chElem, err := ch.NewElement()
		// if err != nil {
		// 	log.Println(err)
		// }

		// fmt.Println(chElem)
	})
}

type Simple struct {
	A string
	B bool
	C []int    `struct:"true"`
	D []string `struct:"true"`
	E []bool
}

var SS = []Simple{
	{"1", true, []int{0, 1, 2, 3}, []string{"a", "b"}, []bool{true, false, false}},
	{"2", true, []int{40, 41}, []string{"c", "d", "e"}, []bool{true, false, false}},
	{"3", true, []int{51, 52, 53}, []string{"f"}, []bool{true, false, false}},
}

func TestZ(t *testing.T) {
	t.Run("Z", func(t *testing.T) {
		s := meta.ToStructs(SS)
		children := s[0].Fields().WithTagTrue("struct")
		data := s.ExtractDataByName(children.Names()...)
		for k, v := range data {
			fmt.Printf("%s\t%v\n", k, v)
			for _, e := range v {
				// fmt.Println(e)
				fmt.Printf("%T\n", e)
			}
		}
		// fmt.Printf("%#v\n", data)
	})
}
