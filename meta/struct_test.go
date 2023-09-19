package meta_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/exiledavatar/gotoolkit/meta"
)

func TestToStruct(t *testing.T) {

	var testCases = []struct {
		Name  string
		Input any
		// Expect meta.Fields
	}{
		{Name: "struct: ExampleStruct", Input: structExample},
		{Name: "struct: empty", Input: struct{}{}},
		{Name: "struct: &ExampleStruct", Input: &structExample},
		{Name: "nil", Input: nil},
		{Name: "reflect.Value: ExampleStruct", Input: reflect.ValueOf(structExample)},
		{Name: "reflect.Value: &ExampleStruct", Input: reflect.ValueOf(&structExample)},
		{Name: "reflect.Value: nil", Input: reflect.ValueOf(nil)},
		{Name: "reflect.Type: ExampleStruct", Input: reflect.TypeOf(structExample)},
		{Name: "reflect.Type: &ExampleStruct", Input: reflect.TypeOf(&structExample)},
		{Name: "reflect.Type: nil", Input: reflect.TypeOf(nil)},
	}
	fmt.Printf("%-40s%-30s%-30s%-30s%-12s\n", "Test", "Input", "Output(Length)", "Expect(Length)", "Pass")
	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			output, err := meta.ToStruct(v.Input)
			// fieldsSliceSlice = append(fieldsSliceSlice, fields)
			// pass := len(fields) == len(v.Expect)

			fmt.Printf("%-40s%-30T%-30s%-50v\n",
				v.Name,
				v.Input,
				fmt.Sprintf("%T", output),
				err,
			)
			// fmt.Println(output.Childs().Names())
			// fmt.Sprintf("%T(%d)", v.Expect, len(v.Expect)),
		})
	}

}

func TestStruct_ExecuteTemplate(t *testing.T) {

	t.Run("Struct.ExecuteTemplate", func(t *testing.T) {
		fmt.Printf("%-.200s\n", fmt.Sprintf("Test: Struct.ExecuteTemplate %-s", strings.Repeat("-", 200)))
		s, err := meta.ToStruct(structExample)
		if err != nil {
			log.Println(err)
		}
		tpl, err := s.ExecuteTemplate(`		
join: {{ join ", " .Struct.Fields.Names }}
{{ .Struct.Fields.Names }}
{{ $fieldsTypes := .Struct.Fields.Types | tostrings }}
joinslices: 
	{{ joinslices ": " "\n\t" .Struct.Fields.Names $fieldsTypes }}
		`, nil, map[string]any{"test": "me"})
		if err != nil {
			log.Println(err)
		}
		fmt.Println(tpl)

	})
	fmt.Printf("%-.200s\n", strings.Repeat("-", 200))

}
