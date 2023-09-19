package meta_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/exiledavatar/gotoolkit/meta"
)

func TestToFields(t *testing.T) {
	// this is the standard output that all non-nil returns should produce
	expectedFields, _ := meta.ToFields(structExample) //

	var nilFields meta.Fields

	var testCases = []struct {
		Name   string
		Input  any
		Expect meta.Fields
	}{
		{Name: "struct: ExampleStruct", Input: structExample, Expect: expectedFields},
		{Name: "struct: empty", Input: struct{}{}, Expect: nilFields},
		{Name: "struct: &ExampleStruct", Input: &structExample, Expect: expectedFields},
		{Name: "nil", Input: nil, Expect: nilFields},
		{Name: "reflect.Value: ExampleStruct", Input: reflect.ValueOf(structExample), Expect: expectedFields},
		{Name: "reflect.Value: &ExampleStruct", Input: reflect.ValueOf(&structExample), Expect: expectedFields},
		{Name: "reflect.Value: nil", Input: reflect.ValueOf(nil), Expect: nilFields},
		{Name: "reflect.Type: ExampleStruct", Input: reflect.TypeOf(structExample), Expect: expectedFields},
		{Name: "reflect.Type: &ExampleStruct", Input: reflect.TypeOf(&structExample), Expect: expectedFields},
		{Name: "reflect.Type: nil", Input: reflect.TypeOf(nil), Expect: nilFields},
	}

	fmt.Printf("%-40s%-30s%-30s%-30s%-12s\n", "Test", "Input", "Output(Length)", "Expect(Length)", "Pass")
	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			fields, _ := meta.ToFields(v.Input)
			// fieldsSliceSlice = append(fieldsSliceSlice, fields)
			pass := len(fields) == len(v.Expect)

			fmt.Printf("%-40s%-30T%-30s%-30s%-12t\n",
				v.Name,
				v.Input,
				fmt.Sprintf("%T(%d)", fields, len(fields)),
				fmt.Sprintf("%T(%d)", v.Expect, len(v.Expect)),
				pass)
		})
	}
}
