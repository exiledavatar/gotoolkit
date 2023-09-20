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
	Name               string   // default is type (without package namespace), single value children structs default to field name, slices default to type, maps default to key name
	NameSpace          []string // default is []string{"package"}
	NameSpaceSeparator string   // defaults to "."
	Value
	UUID       string
	Attributes map[string]any
	Tags       Tags //map[string][]string
	Parent     *Struct
	Data
}

func ToStruct(value any) (Struct, error) {
	d := ToData(value)
	s, isStruct := value.(Struct)
	if isStruct {
		// log.Printf("%v is already a struct\n", value)
		return s, nil
	}
	v, err := ToValue(value)
	switch kind := v.Value.Kind(); {
	case err != nil || kind == reflect.Invalid:
		return s, err
	case (kind == reflect.Slice || kind == reflect.Array) && v.Len() > 0:
		v, err = ToValue(v.Value.Index(0).Interface())
		if err != nil {
			return s, err
		}
	case (kind == reflect.Slice || kind == reflect.Array) && v.Len() == 0:
		v, err = ToValue(reflect.New(v.Type().Elem()).Elem().Interface())
		if err != nil {
			return s, err
		}

	}

	switch kind := v.Kind(); {
	case kind == reflect.Struct:
		pkgPath := strings.Split(v.Type().PkgPath(), "/")
		s = Struct{
			Name:               v.Type().Name(),
			NameSpace:          []string{pkgPath[len(pkgPath)-1]},
			NameSpaceSeparator: ".",
			Value:              v,
			UUID:               strings.ReplaceAll(uuid.NewString(), "-", ""),
			Data:               d,
		}
		return s, nil
	default:
		return s, fmt.Errorf("invalid (kind) type: (%s) %s", v.Kind(), v.Type())
	}
}

// NewUUID creates a new UUID and sets it recursively
func (s *Struct) NewUUID() string {
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	s.SetUUID(id)
	return id
}

// SetUUID recursively sets s and its fields' struct UUID's to id
func (s *Struct) SetUUID(id string) {
	s.UUID = id
}

type Structconfig struct {
	Name                     string
	NameSpace                []string
	NameSpaceSeparator       string // default is "."
	UUID                     string
	Tags                     Tags           // struct tags, not field tags, may optionally be parsed from tags labeled struct
	RemoveExistingTags       bool           // remove existing tags - false will simply apopend any included tags to current ones
	Attributes               map[string]any //string // these should be table constraints, not field constraints
	RemoveExistingAttributes bool           // remove existing constraints - false will simply apopend any included constraints to current ones
	Parent                   *Struct
	// Tags                     map[string][]string // struct tags, not field tags, may optionally be parsed from tags labeled struct
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

	if cfg.Parent != nil {
		ds.Parent = cfg.Parent
	}
	return ds, nil
}

// TagName ranges through the provided keys in order and returns the
// first non-blank, non-false value, or Struct.Name if none are found.
func (s Struct) TagName(keys ...string) string {
	name := s.Name
	for _, key := range keys {
		switch tag := s.Tags.Tag(key); {
		case tag == nil || tag.False():
			continue
		case tag[0] != "":
			return tag[0]
		default:
			// return f.Name
		}

	}
	return name
}

// Identifier returns the full namespaced identifier of the struct
func (s Struct) Identifier() string {
	ids := append(s.NameSpace, s.Name)
	for i, v := range ids {
		ids[i] = strings.ToLower(v)
	}
	return strings.Join(ids, s.NameSpaceSeparator)
}

// Identifier returns the full namespaced identifier of the struct,
// but uses TagName instead of Name
func (s Struct) TagIdentifier(keys ...string) string {
	ids := append(s.NameSpace, s.TagName(keys...))
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
func (s *Struct) ExecuteTemplate(tpl string, funcs template.FuncMap, data map[string]any) (string, error) {
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
	return fields
}

func (s *Struct) ExtractDataByName(names ...string) map[string]Data {
	data := map[string]Data{}
	data[s.Name] = Data(ToSlice(s.Value.Interface()))
	for _, child := range s.Fields().ByNames(names...) {
		data[child.Name] = Data(ToSlice(child.Value.Interface()))
	}
	return data
}

func (s *Struct) Extract(names ...string) map[string]Struct {
	structs := map[string]Struct{}
	for _, child := range s.Fields().ByNames(names...) {
		childStruct, err := child.ToStruct()
		if err != nil {
			continue
		}
		structs[child.Name] = childStruct
	}

	data := map[string]Data{}
	for _, row := range s.Data {
		rowValue := reflect.ValueOf(row)
		for _, childName := range names {
			childRowData := ToData(rowValue.FieldByName(childName).Interface())
			data[childName] = append(data[childName], childRowData...)
		}
	}
	for k, str := range structs {
		str.Data = data[k]
		structs[k] = str
	}
	structs[s.Name] = *s
	return structs
}
