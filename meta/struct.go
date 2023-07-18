package meta

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

// Struct captures the properties of a struct and allows overriding properties using tags or by direct assignment
type Struct struct {
	Name               string // default is type, single value children structs default to field name, slices default to type, maps default to key name
	NameSpace          []string
	NameSpaceSeparator string // defaults to "."
	Value
	UUID       string
	Attributes map[string]string
	Tags       map[string][]string
	Parent     *Struct
}

func ToStruct(value any) (Struct, error) {

	s, isStruct := value.(Struct)
	if isStruct {
		log.Printf("%v is already a struct\n", value)
		return s, nil
	}

	switch v, err := ToValue(value); {
	case err != nil:
		log.Println("ToValue err not nil")
		return s, err
	case v.Kind() != reflect.Struct:
		log.Println("ToValue(value).Kind() != reflect.Struct")
		return s, fmt.Errorf("invalid (kind) type: (%s) %s", v.Kind(), v.Type())
	default:
		log.Println("default, creating struct")
		s = Struct{
			Name:               v.Name,
			NameSpaceSeparator: ".",
			Value:              v,
			UUID:               strings.ReplaceAll(uuid.NewString(), "-", ""),
		}
		return s, nil
	}
}

// func (s Struct) Childs() []Struct {
// 	var ss []Struct

// 	// for i, field := range s.Fields {
// 	// 	if len(field.Struct.Fields) > 0 {
// 	// 		ss = append(ss, s.Fields[i].Struct)
// 	// 	}
// 	// }
// 	return ss
// }

// NewUUID creates a new UUID and sets it recursively
func (s *Struct) NewUUID() string {
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	s.SetUUID(id)
	return id
}

// SetUUID recursively sets s and its fields' struct UUID's to id
func (s *Struct) SetUUID(id string) {
	s.UUID = id
	// if s.Fields != nil {
	// 	s.Fields.SetUUID(id)
	// }
}

type Structconfig struct {
	Name                     string
	NameSpace                []string
	NameSpaceSeparator       string // default is "."
	UUID                     string
	Tags                     map[string][]string // struct tags, not field tags, may optionally be parsed from tags labeled struct
	RemoveExistingTags       bool                // remove existing tags - false will simply apopend any included tags to current ones
	Attributes               map[string]string   // these should be table constraints, not field constraints
	RemoveExistingAttributes bool                // remove existing constraints - false will simply apopend any included constraints to current ones
}

// NewStruct enables configuration when parsing a struct to a Struct
func NewStruct(value any, cfg Structconfig) (Struct, error) {
	ds, err := ToStruct(value)
	if err != nil {
		return ds, err
	}
	if cfg.Name != "" {
		ds.Name = cfg.Name
	}

	if len(cfg.NameSpace) > 0 {
		ds.NameSpace = cfg.NameSpace
	}
	if cfg.UUID != "" {
		ds.UUID = cfg.UUID
	}

	if cfg.RemoveExistingTags {
		ds.Tags = nil
	}

	switch {
	case ds.Tags == nil && cfg.Tags != nil:
		ds.Tags = cfg.Tags
	case ds.Tags != nil && cfg.Tags != nil:
		for k, v := range cfg.Tags {
			ds.Tags[k] = v
		}
	}

	if cfg.RemoveExistingAttributes {
		ds.Attributes = nil
	}

	switch {
	case ds.Attributes == nil && cfg.Attributes != nil:
		ds.Attributes = cfg.Attributes
	case ds.Attributes != nil && cfg.Attributes != nil:
		for k, v := range cfg.Attributes {
			ds.Attributes[k] = v
		}
	}

	return ds, nil
}

func (s Struct) Identifier() string {
	ids := append(s.NameSpace, s.Name)
	for i, v := range ids {
		ids[i] = strings.ToLower(v)
	}
	return strings.Join(ids, s.NameSpaceSeparator)
}

func (s Struct) ValueMap(tagKey string) ValueMap {
	return ToValueMap(s, tagKey)
}

// ExecuteTemplate parses a string and executes it with any additional funcs and data. All data, including the reciever
// is passed to text/template as a map. By default, the reciever's map key is its type - eg {{ .Struct }} references a calling Struct.
// By default, it passes missingkey=zero, you can override this by changing TemplateOptions
// See TemplateFuncMap for additional functions included by default.
// See TemplateDataNames if you really need to change data map key names.
func (s *Struct) ExecuteTemplate(tpl string, funcs template.FuncMap) (string, error) {
	// d := map[string]any{
	// 	TemplateDataNames["Struct"]: s,
	// }
	// for k, v := range data {
	// 	d[k] = v
	// }
	data := s
	parsedTpl, err := template.
		New("").
		Option(TemplateOptions...).
		Funcs(TemplateFuncMap).
		Funcs(funcs).
		Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := parsedTpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// // ExecuteTemplate parses a string and executes it with any additional funcs and data. All data, including the reciever
// // is passed to text/template as a map. By default, the reciever's map key is its type - eg {{ .Struct }} references a calling Struct.
// // By default, it passes missingkey=zero, you can override this by changing TemplateOptions
// // See TemplateFuncMap for additional functions included by default.
// // See TemplateDataNames if you really need to change data map key names.
// func (s Struct) ExecuteTemplate(tpl string, funcs template.FuncMap, data map[string]any) (string, error) {
// 	d := map[string]any{
// 		TemplateDataNames["Struct"]: s,
// 	}
// 	for k, v := range data {
// 		d[k] = v
// 	}

// 	parsedTpl, err := template.
// 		New("").
// 		Option(TemplateOptions...).
// 		Funcs(TemplateFuncMap).
// 		Funcs(funcs).
// 		Parse(tpl)
// 	if err != nil {
// 		return "", err
// 	}

// 	var buf bytes.Buffer
// 	if err := parsedTpl.Execute(&buf, d); err != nil {
// 		return "", err
// 	}

// 	return buf.String(), nil
// }

func (s *Struct) Fields() Fields {

	sfs := reflect.VisibleFields(s.Value.Type())
	sfmap := map[string]reflect.StructField{}
	for _, sf := range sfs {
		if sf.IsExported() && !sf.Anonymous {
			sfmap[sf.Name] = sf
		}
	}

	var fields Fields
	for _, child := range s.Value.Children() {
		field := Field{
			Name:        child.Name,
			Parent:      s,
			Value:       child,
			StructField: sfmap[child.Name],
		}
		fields = append(fields, field)
	}
	// fields, err := ToFields(s.Children())
	// if err != nil {
	// 	panic(err)
	// }
	// for i := range fields {
	// 	fields[i].Parent = s
	// }
	return fields
}

func (s Struct) Fields2() Fields {

	sfs := reflect.VisibleFields(s.Value.Type())
	sfmap := map[string]reflect.StructField{}
	for _, sf := range sfs {
		if sf.IsExported() && !sf.Anonymous {
			sfmap[sf.Name] = sf
		}
	}

	var fields Fields
	for _, child := range s.Value.Children() {
		field := Field{
			Name:        child.Name,
			Parent:      &s,
			Value:       child,
			StructField: sfmap[child.Name],
		}
		fields = append(fields, field)
	}
	// fields, err := ToFields(s.Children())
	// if err != nil {
	// 	panic(err)
	// }
	// for i := range fields {
	// 	fields[i].Parent = s
	// }
	return fields
}

func (s *Struct) Fields3() Fields {

	sfs := reflect.VisibleFields(s.Value.Type())
	sfmap := map[string]reflect.StructField{}
	for _, sf := range sfs {
		if sf.IsExported() && !sf.Anonymous {
			sfmap[sf.Name] = sf
		}
	}

	var fields Fields
	for _, child := range s.Value.Children() {
		field := Field{
			Name:        child.Name,
			Parent:      s,
			Value:       child,
			StructField: sfmap[child.Name],
		}
		fields = append(fields, field)
	}
	// fields, err := ToFields(s.Children())
	// if err != nil {
	// 	panic(err)
	// }
	// for i := range fields {
	// 	fields[i].Parent = s
	// }
	return fields
}

// func (s Struct) ExecuteTemplate(tpl string, funcs template.FuncMap) string {
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New(s.Name).
// 		Option("missingkey=zero").
// 		Funcs(FuncMap).
// 		Funcs(funcs).
// 		Parse(tpl)).
// 		Execute(&buf, s); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// }

// func (s Struct) ExecuteTemplateWithData(tpl string, funcs template.FuncMap, data ...any) string {
// 	data = append([]any{s}, data...)
// 	var buf bytes.Buffer

// 	if err := template.Must(template.
// 		New("").
// 		Option("missingkey=zero").
// 		Funcs(FuncMap).
// 		Funcs(funcs).
// 		Parse(tpl)).
// 		Execute(&buf, data); err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	return buf.String()
// }
