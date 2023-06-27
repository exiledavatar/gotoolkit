package meta

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

// Struct captures the properties of a struct and allows overriding properties using tags or by direct assignment
type Struct struct {
	Name               string
	NameSpace          []string
	NameSpaceSeparator string        // defaults to "."
	Type               reflect.Type  // pointers will be de-referenced (indirected)
	Value              reflect.Value // pointers will be de-referenced (indirected)
	UUID               string
	Attributes         map[string]string
	Fields             Fields
	Children           []Structs // fields that are slices of structs with more than one member
	Tags               map[string][]string
	Parent             *Struct
	Container          reflect.Type // map, slice, or array Struct was 'wrapped' in, if applicable
	// Data               any
	pointer bool // is the original variable a pointer
}

func ToStruct(value any) (Struct, error) {
	var s Struct
	rv, rt, pointer := ToIndirectReflectValue(value)
	if rt == nil {
		return s, fmt.Errorf("invalid value: nil")

	}
	if rt.Kind() != reflect.Struct {
		return s, fmt.Errorf("invalid type: (%s) %s", rt.Kind(), rt)
	}
	switch {
	case rt == nil:
		return s, fmt.Errorf("invalid value: %v", value)
	case rt.Kind() == reflect.Invalid:
		return s, fmt.Errorf("invalid value: Kind() == reflect.Invalid: %v", value)
	case rt.Kind() == reflect.Map:

	case rt.Kind() == reflect.Slice:
	case rt.Kind() == reflect.Array:
	case rt.Kind() == reflect.Chan:
	case rt.Kind() != reflect.Struct:
		return s, fmt.Errorf("invalid type: (%s) %s", rt.Kind(), rt)
	}
	s = Struct{
		Name:       rt.Name(),
		Type:       rt,
		Value:      rv,
		Attributes: nil,
		Fields:     ToFields(value),
		// Data:       value,
		pointer: pointer,
		// UUID:       strings.ReplaceAll(uuid.NewString(), "-", ""),
	}
	s.NewUUID()

	for i, field := range s.Fields {
		if len(field.Struct.Fields) > 0 {
			s.Fields[i].Struct.Parent = &s
		}
	}

	return s, nil
}

func (s Struct) Childs() []Struct {
	var ss []Struct

	for i, field := range s.Fields {
		if len(field.Struct.Fields) > 0 {
			ss = append(ss, s.Fields[i].Struct)
		}
	}
	return ss
}

// NewUUID creates a new UUID and set it recursively
func (s *Struct) NewUUID() string {
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	s.SetUUID(id)
	return id
}

// SetUUID recursively sets s and its fields' struct UUID's to id
func (s *Struct) SetUUID(id string) {
	s.UUID = id
	if s.Fields != nil {
		s.Fields.SetUUID(id)
	}
}

type StructConfig struct {
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
func NewStruct(a any, cfg StructConfig) (Struct, error) {
	ds, err := ToStruct(a)
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
func (s Struct) ExecuteTemplate(tpl string, funcs template.FuncMap, data map[string]any) (string, error) {
	d := map[string]any{
		TemplateDataNames["Struct"]: s,
	}
	for k, v := range data {
		d[k] = v
	}

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
	if err := parsedTpl.Execute(&buf, d); err != nil {
		return "", err
	}

	return buf.String(), nil
}
