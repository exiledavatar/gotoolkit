package meta_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/ericsgagnon/qgenda/pkg/meta"
)

func TestToValueMap(t *testing.T) {

	var testCases = []struct {
		Name  string
		Input any
		// Expect meta.Fields
	}{
		{Name: "struct: ExampleStruct", Input: structExample},
		{Name: "any(struct): ExampleStruct", Input: any(structExample)},
		{Name: "struct: &ExampleStruct", Input: &structExample},
		{Name: "struct: empty", Input: struct{}{}},
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
			output := meta.ToValueMap(v.Input, "vm")
			b, err := json.Marshal(output)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%-40s%-30T%+.200v\n",
				v.Name,
				v.Input,
				string(b),
			)

		})
	}

}

// func TestToValue(t *testing.T) {

// }
