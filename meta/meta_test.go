package meta_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/exiledavatar/gotoolkit/meta"
)

type Struct0Fields struct{}

type Struct1Fields struct {
	One int
}

type Struct2Fields struct {
	One int
	Two int
}

type ExStruct struct {
	Bool        bool
	BoolPointer *bool
}

type ExampleStruct struct {
	BoolValue           bool             `vm:"true" json:"boolvalue"`
	boolValue           bool             `vm:""`
	BoolPointer         *bool            `vm:"" json:"boolpointer,omitempty"`
	boolPointer         *bool            `vm:""`
	StringValue         string           `vm:"" db:"stringvalue"`
	stringValue         string           `vm:""`
	StringPointer       *string          `vm:"" db:"dbstringpointer" pg:"pgstringpointer" json:"jsonstringpointer"`
	stringPointer       *string          `vm:""`
	IntValue            int              `vm:"" db:"dbintvalue" pg:"pgintvalue"`
	intValue            int              `vm:""`
	IntPointer          *int             `vm:"" pg:"pgintpointer"`
	intPointer          *int             `vm:""`
	Bytes               []byte           `vm:""`
	Map                 map[string]any   `vm:""`
	Slice               []any            `vm:"" db:"dbslice" pg:""`
	slice               *[]any           `vm:""`
	StructValue         ExStruct         `vm:""`
	structValue         ExStruct         `vm:""`
	StructPointer       *ExStruct        `vm:""`
	structPointer       *ExStruct        `vm:""`
	ExampleStructSlice  []ExampleStruct  `vm:"" struct:"true"`
	exampleStructSlice  *[]ExampleStruct `vm:""`
	ExampleStructSlice2 []ExampleStruct  `vm:"" struct:"true"`
	ExampleStructSlice3 []ExampleStruct  `vm:"" struct:"true"`
	Struct0Fields       Struct0Fields
	Struct1Fields       Struct1Fields
	Struct2Fields       Struct2Fields
	ExStruct
}

var boolExample = true
var stringExample = "example"
var intExample = 1138
var byteExample byte = 42
var byteSliceExample = []byte{7, 8, 42}
var mapExample = map[string]any{
	"bool":          boolExample,
	"boolPointer":   &boolExample,
	"string":        stringExample,
	"stringPointer": &stringExample,
	"int":           intExample,
	"intPointer":    &intExample,
	"struct":        ExampleStruct{},
	"nil":           nil,
}

var sliceExample = []any{
	boolExample,
	&boolExample,
	stringExample,
	&stringExample,
	intExample,
	&intExample,
	byteExample,
	&byteExample,
	nil,
}

var structExample = ExampleStruct{
	BoolValue:          boolExample,
	boolValue:          boolExample,
	BoolPointer:        &boolExample,
	boolPointer:        &boolExample,
	StringValue:        stringExample,
	stringValue:        stringExample,
	StringPointer:      &stringExample,
	stringPointer:      &stringExample,
	IntValue:           intExample,
	intValue:           intExample,
	IntPointer:         &intExample,
	intPointer:         &intExample,
	Bytes:              byteSliceExample,
	Map:                mapExample,
	Slice:              sliceExample,
	slice:              &sliceExample,
	StructValue:        ExStruct{},
	structValue:        ExStruct{},
	StructPointer:      &ExStruct{},
	structPointer:      &ExStruct{},
	ExampleStructSlice: []ExampleStruct{},
	exampleStructSlice: &[]ExampleStruct{},
	ExStruct:           ExStruct{},
}

func TestToIndirectReflectValue(t *testing.T) {
	// this is the standard output that all non-nil returns should produce
	expectedRV, expectedRT, _ := meta.ToIndirectReflectValue(structExample) //

	type Expect struct {
		Value   reflect.Value
		Type    reflect.Type
		Pointer bool
	}
	var testCases = []struct {
		Name   string
		Input  any
		Expect Expect
	}{
		{Name: "struct: ExampleStruct", Input: structExample, Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: false}},
		{Name: "struct: empty", Input: struct{}{}, Expect: Expect{Value: reflect.ValueOf(struct{}{}), Type: reflect.TypeOf(struct{}{}), Pointer: false}},
		{Name: "struct: &ExampleStruct", Input: &structExample, Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: true}},
		{Name: "nil", Input: nil, Expect: Expect{Value: reflect.ValueOf(nil), Type: reflect.TypeOf(nil), Pointer: false}},
		{Name: "reflect.Value: ExampleStruct", Input: reflect.ValueOf(structExample), Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: false}},
		{Name: "reflect.Value: &ExampleStruct", Input: reflect.ValueOf(&structExample), Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: true}},
		{Name: "reflect.Value: nil", Input: reflect.ValueOf(nil), Expect: Expect{Value: reflect.ValueOf(nil), Type: reflect.TypeOf(nil), Pointer: false}},
		{Name: "reflect.Type: ExampleStruct", Input: reflect.TypeOf(structExample), Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: false}},
		{Name: "reflect.Type: &ExampleStruct", Input: reflect.TypeOf(&structExample), Expect: Expect{Value: expectedRV, Type: expectedRT, Pointer: true}},
		{Name: "reflect.Type: nil", Input: reflect.TypeOf(nil), Expect: Expect{Value: reflect.ValueOf(nil), Type: reflect.TypeOf(nil), Pointer: false}},
		{Name: "map[string]any", Input: mapExample, Expect: Expect{Value: reflect.ValueOf(mapExample), Type: reflect.TypeOf(mapExample), Pointer: false}},
	}
	fmt.Printf("%-.200s\n", fmt.Sprintf("Test: ToIndirectReflectValue %-s", strings.Repeat("-", 200)))
	// fmt.Println(fmt.Sprintf("%-50s", "Test: ToIndirectReflectValue"))
	fmt.Printf("%-40s%-30s%-12s%-40s%-8s%-12s%-40s%-8s%-12s\n", "Test", "Input", "OutputKind", "OutputType", "Pointer", "ExpectKind", "ExpectType", "Pointer", "Pass")
	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			rv, rt, pointer := meta.ToIndirectReflectValue(v.Input)
			// fields := meta.ToFields(v.Input)
			// fieldsSliceSlice = append(fieldsSliceSlice, fields)
			// pass := rv.Equal(v.Expect.Value) && pointer == v.Expect.Pointer
			// pass := reflect.DeepEqual(rv, v.Expect.Value)
			pass := rt == v.Expect.Type && pointer == v.Expect.Pointer //&& rt.Kind() == v.Expect.Value.Kind()
			fmt.Printf("%-40s%-30T%-60s%-60s%-12t\n",
				v.Name,
				v.Input,
				fmt.Sprintf("%-12s%-40s%-8t", rv.Kind(), fmt.Sprint(rt), pointer),
				fmt.Sprintf("%-12s%-40s%-8t", v.Expect.Value.Kind(), fmt.Sprint(v.Expect.Type), v.Expect.Pointer),
				pass)
		})
	}
	// colorReset := "\033[0m"
	// colorRed := "\033[31m"
	// // colorGreen := "\033[32m"
	// fmt.Println(string(colorRed), "test", string(colorReset))
	// fmt.Println(colorRed, "\033[32mnext")

}

func TestToSlicey(t *testing.T) {
	cases := []any{
		[]string{"a", "b", "c"},
		[]int{1, 2, 3},
		"zed",
		1432,
		[][]string{{"98", "76"}},
		[]any{[]any{1, 4, 2}, []string{"h", "i", "j"}},
	}
	for i, v := range cases {
		vs := meta.ToSlice(v)
		fmt.Printf("%d\t%#v\t\t%#v\n", i, vs, meta.Flatten(vs))

	}

	fmt.Printf("%#v\n", meta.ToSlice(cases[1], cases[3], ""))
	// panic(0)
}
