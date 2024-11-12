package meta

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"reflect"
)

// TemplateFuncMap includes functions in addition to those already included
// by text/template: https://pkg.go.dev/text/template#hdr-Functions
var TemplateFuncMap = template.FuncMap{
	"coalesce":      Coalesce[[]string, string],
	"join":          join, // wraps strings.Join but reorders arguments to work with template piping
	"joinslices":    JoinSlices,
	"mapkeys":       mapKeys,
	"mapvalues":     mapValues,
	"replace":       replace,
	"replaceall":    replaceAll,
	"tolower":       strings.ToLower,
	"tolowerslices": toLowerSlices,
	"toupper":       strings.ToUpper,
	"toupperslices": toUpperSlices,
	"tostrings":     toStrings,
	"toslice":       toSlice,
	"tovaluemap":    ToValueMap,
}

func replace(old string, new string, n int, s string) string {
	return strings.Replace(s, old, new, n)
}

func replaceAll(old, new, s string) string {
	return strings.ReplaceAll(s, old, new)
}

// toLowerSlices takes a mix of string and []string and returns them in a flattened, lowercased []string
func toLowerSlices(values ...any) []string {
	var out []string
	for _, value := range ToStringSlice(values...) {
		out = append(out, strings.ToLower(value))
	}
	return out
}

// toUpperSlices takes a mix of string and []string and returns them in a flattened, uppercased []string
func toUpperSlices(values ...any) []string {
	var out []string
	for _, value := range ToStringSlice(values...) {
		out = append(out, strings.ToUpper(value))
	}
	return out
}

// TODO: replace toSlice with ToSlice
func toStrings(s any) []string {
	strings := []string{}
	slicey := toSlice(s)
	for _, v := range slicey {
		strings = append(strings, fmt.Sprint(v))
	}
	return strings
}

// TODO: this should be deprecated in favor of ToSlice
// toSlice is intended to explicity convert a slice to a slice of type any
func toSlice(a any) []any {
	v := reflect.ValueOf(a)
	var s []any
	if v.Kind() != reflect.Slice {
		log.Printf("%T is not a slice\n", a)
	}
	if v.Kind() == reflect.Slice {
		iv := reflect.Indirect(v)
		sliceType := reflect.TypeOf(a).Elem()
		out := reflect.MakeSlice(reflect.SliceOf(sliceType), iv.Len(), iv.Len())
		for i := 0; i < iv.Len(); i++ {
			f := reflect.Indirect(iv.Index(i))
			out.Index(i).Set(f)
			s = append(s, f.Interface())
		}
	}
	return s
}

// join wraps strings.Join to enable pipelining in templates
func join(sep string, s []string) string {
	return strings.Join(s, sep)
}

func JoinSlices(elementSep string, indexSep string, slices ...[]string) string {
	// check for equal length slices
	var lengths []int
	ok := true
	for i, sl := range slices {
		lengths = append(lengths, len(sl))
		ok = lengths[0] > 0 && lengths[i] == lengths[0]
	}
	if !ok {
		panic(fmt.Sprint("length of each slice must be the same, actual lengths are ", lengths))
	}

	elements := []string{}
	for i := 0; i < lengths[0]; i++ {
		ielements := []string{}
		for j := 0; j < len(slices); j++ {
			ielements = append(ielements, fmt.Sprint(slices[j][i]))
		}
		elements = append(elements, strings.Join(ielements, elementSep))
	}
	return strings.Join(elements, indexSep)
}

func mapKeys(m map[string]any) []string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, fmt.Sprint(k))
	}
	return keys
}

func mapValues(m map[string]any) []string {
	values := []string{}
	for _, v := range m {
		values = append(values, fmt.Sprint(v))
	}
	return values
}

var TemplateDataNames = map[string]string{
	"Struct": "Struct",
	// "Structs": "Structs",
	// "Field":   "Field",
	// "Fields":  "Fields",
	// "Tag":     "Tag",
	// "Tags":    "Tags",
}

// TemplateOptions is passed to all functions/methods that execute a template.Template
// see https://pkg.go.dev/text/template#Template.Option
var TemplateOptions = []string{
	"missingkey=zero",
}
