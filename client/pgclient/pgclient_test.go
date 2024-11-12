package pgclient_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/exiledavatar/gotoolkit/client/pgclient"
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
	IDHash              string           `pg:"id_hash" primarykey:"true"`
	BoolValue           bool             `vm:"true" json:"boolvalue"`
	boolValue           bool             `vm:""`
	BoolPointer         *bool            `vm:"" json:"boolpointer,omitempty"`
	boolPointer         *bool            `vm:""`
	StringValue         string           `vm:"" db:"stringvalue" pgtype:"text"`
	stringValue         string           `vm:""`
	StringPointer       *string          `vm:"" db:"dbstringpointer" pg:"pgstringpointer" json:"jsonstringpointer"`
	stringPointer       *string          `vm:""`
	IntValue            int              `vm:"" db:"dbintvalue" pg:"pgintvalue" pgtype:"bigint"`
	intValue            int              `vm:""`
	IntPointer          *int             `vm:"" pg:"pgintpointer" pgtype:"something_silly"`
	intPointer          *int             `vm:""`
	Time                time.Time        `pg:"some_time"`
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

func Test(t *testing.T) {
	// this is the standard output that all non-nil returns should produce

	// expectedFields, _ := meta.ToFields(structExample) //

	// var nilFields meta.Fields

	// var testCases = []struct {
	// 	Name   string
	// 	Input  any
	// 	Expect meta.Fields
	// }{
	// 	{Name: "struct: ExampleStruct", Input: structExample, Expect: expectedFields},
	// 	{Name: "struct: empty", Input: struct{}{}, Expect: nilFields},
	// 	{Name: "struct: &ExampleStruct", Input: &structExample, Expect: expectedFields},
	// 	{Name: "nil", Input: nil, Expect: nilFields},
	// 	{Name: "reflect.Value: ExampleStruct", Input: reflect.ValueOf(structExample), Expect: expectedFields},
	// 	{Name: "reflect.Value: &ExampleStruct", Input: reflect.ValueOf(&structExample), Expect: expectedFields},
	// 	{Name: "reflect.Value: nil", Input: reflect.ValueOf(nil), Expect: nilFields},
	// 	{Name: "reflect.Type: ExampleStruct", Input: reflect.TypeOf(structExample), Expect: expectedFields},
	// 	{Name: "reflect.Type: &ExampleStruct", Input: reflect.TypeOf(&structExample), Expect: expectedFields},
	// 	{Name: "reflect.Type: nil", Input: reflect.TypeOf(nil), Expect: nilFields},
	// }

	// fmt.Printf("%-40s%-30s%-30s%-30s%-12s\n", "Test", "Input", "Output(Length)", "Expect(Length)", "Pass")
	// for _, v := range testCases {
	// 	t.Run(v.Name, func(t *testing.T) {
	// 		fields, _ := meta.ToFields(v.Input)
	// 		// fieldsSliceSlice = append(fieldsSliceSlice, fields)
	// 		pass := len(fields) == len(v.Expect)

	// 		fmt.Printf("%-40s%-30T%-30s%-30s%-12t\n",
	// 			v.Name,
	// 			v.Input,
	// 			fmt.Sprintf("%T(%d)", fields, len(fields)),
	// 			fmt.Sprintf("%T(%d)", v.Expect, len(v.Expect)),
	// 			pass)
	// 	})
	// }

	// str, err := meta.ToStruct(structExample)
	// if err != nil {
	// 	log.Println(err)
	// }
	// out, err := str.ExecuteTemplate(
	// 	pgclient.PGTemplates.CreateTable,
	// 	nil,
	// 	map[string]any{
	// 		"Config": pgclient.TemplateConfig,
	// 	},
	// )
	// switch {
	// case err != nil:
	// 	log.Println(err)
	// default:
	// 	fmt.Println(out)
	// }
	text, err := pgclient.CreateTableText(structExample)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(text)
}
