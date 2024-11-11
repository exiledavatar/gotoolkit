package meta_test

import (
	"testing"
)

var dataTestCases = []any{
	true,
	[]bool{true, false, false, true, false},
	"string value",
	[]string{"A", "B", "C", "D", "E", "F"},
	1138,
	[]int{100, 99, 98, 97, 96, 95, 94},
	byte(13),
	[]byte{7, 8, 42},
	mapExample,
	[]map[string]any{mapExample, mapExample},
	sliceExample,
	[][]any{sliceExample, sliceExample},
	structExample,
	[]ExampleStruct{structExample, structExample},
	nil,
}

func TestToData(t *testing.T) {
	// all non-nil output should be a slice of any who's static element type matches input
	// for _, input := range dataTestCases {
	// 	t.Run(fmt.Sprintf("%T", input), func(t *testing.T) {
	// 		out := meta.ToData(input)
	// 		// fmt.Println(out)
	// 		if input != nil {
	// 			rvInput := reflect.ValueOf(input)
	// 			var rtInput, rtOut reflect.Type
	// 			if (rvInput.Kind() == reflect.Slice || rvInput.Kind() == reflect.Map || rvInput.Kind() == reflect.Array) && rvInput.Len() > 0 {
	// 				rtInput = rvInput.Type()
	// 			} else {
	// 				rtInput = rvInput.Type()
	// 			}
	// 			rvOut := reflect.ValueOf(out)
	// 			if rvOut.Len() > 0 {
	// 				rtOut = rvOut.Index(0).Type()
	// 			}
	// 			fmt.Printf("%T:\t%s\t%s\n", input, rtInput, rtOut)
	// 		}

	// 	})
	// }
	// os.Exit(0)
}

func TestToSlice(t *testing.T) {
	// all non-nil output should be a slice of any who's static element type matches input

}
