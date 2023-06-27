package meta

import (
	"reflect"
	"strings"
)

// Field is a wrapper for reflect.StructField with some
// additional functionality for tags, templating, etc.
type Field struct {
	Name         string
	ForeignTypes map[string]any
	ClientTypes  map[string]string // intended type for named client - eg "postgres": "numeric"
	Attributes   map[string]string // catchall for additional attributes
	reflect.StructField
	Value   reflect.Value
	Parent  *Struct
	pointer bool
	Struct  Struct
}

// type FieldConfig struct {
// 	Name                      string
// 	ClientTypes               map[string]string
// 	RemoveExistingClientTypes bool
// 	Attributes                map[string]string
// 	RemoveExistingAttributes  bool
// }

func (f Field) IsStruct() bool {
	return len(f.Struct.Fields) > 0
}

// Pointer returns true if its reflect.Kind is a reflect.Pointer
func (f Field) Pointer() bool {
	return f.StructField.Type.Kind() == reflect.Pointer
}

// Kind returns the reflect.Kind. It will first dereference a pointer
func (f Field) Kind() reflect.Kind {
	switch sft := f.StructField.Type; {
	case sft.Kind() == reflect.Pointer:
		return sft.Elem().Kind()
	default:
		return sft.Kind()
	}
}

// Type returns the reflect.Type. It will first dereference a pointer
func (f Field) Type() reflect.Type {
	ft := f.StructField.Type
	if ft.Kind() == reflect.Pointer {
		return ft.Elem()
	}
	return ft
}

// ElementType returns the type of an element of a map or slice
func (f Field) ElementType() reflect.Type {
	ft := f.StructField.Type
	if ft.Kind() == reflect.Slice || ft.Kind() == reflect.Map {
		return ft.Elem()
	}
	return nil
}

// Tags returns the parsed struct field tags
func (f Field) Tags() Tags {
	return ToTags(string(f.StructField.Tag))
}

// HasTag returns true if the tag key exists
func (f Field) HasTag(key string) bool {
	return f.Tags().Exists(key)
}

// HasTagValue returns true if it has both the key and value
func (f Field) HasTagValue(key, value string) bool {
	return f.Tags().Contains(key, value)
}

// HasTagTrue returns true if its tags satisfy Tags.True
func (f Field) HasTagTrue(key string) bool {
	return f.Tags().True(key)
}

// HasTagFalse returns true if its tags satisfy Tags.False
func (f Field) HasTagFalse(key string) bool {
	return f.Tags().False(key)
}

func (f *Field) SetUUID(id string) {
	if f.Struct.Name != "" {
		f.Struct.SetUUID(id)
	}
}

func (f Field) ForeignType(target string) any {
	if ft, ok := f.ForeignTypes[target]; ok {
		return ft
	}
	if ft, ok := ForeignTypes[target]; ok {
		return ft
	}
	return nil
}

// TaggedName returns the first value of the field's tag with the given key
// if falls back to a lowercase version of the field's name
func (f Field) TaggedName(key string) string {
	if f.HasTagTrue(key) {
		name := f.Tags()[key][0]
		if name != "" {
			return name
		}
	}
	return strings.ToLower(f.Name)
}
