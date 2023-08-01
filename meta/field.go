package meta

import (
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

// Field is a wrapper for reflect.StructField with some
// additional functionality for tags, templating, etc.
type Field struct {
	Name       string
	Attributes map[string]string // catchall for additional attributes
	Parent     *Struct
	Value
	reflect.StructField
}

// func (f Field) IsStruct() bool {
// 	return len(f.Struct.Fields) > 0
// }

// Pointer returns true if its reflect.Kind is a reflect.Pointer
func (f Field) Pointer() bool {
	// return f.StructField.Type.Kind() == reflect.Pointer
	return f.Value.Pointer
}

// Kind returns the reflect.Kind. It will first dereference a pointer
func (f Field) Kind() reflect.Kind {
	return f.Value.Kind()
}

// Type returns the reflect.Type. It will first dereference a pointer
func (f Field) Type() reflect.Type {
	switch ft := f.StructField.Type; {
	case ft == nil:
		return reflect.TypeOf(nil)
	case ft.Kind() == reflect.Pointer:
		return ft.Elem()
	default:
		return ft
	}
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

// TagName returns the first tag value if the tag key exists and is not blank,
// otherwise it returns Field.Name
func (f Field) TagName(key string) string {
	switch tag := f.Tags().Tag(key); {
	case tag == nil:
		return f.Name
	case tag.False():
		return f.Name
	case tag[0] != "":
		return tag[0]
	default:
		return f.Name
	}
}

func (f Field) Tag(key string) Tag {
	return f.Tags().Tag(key)
}

// func (f *Field) SetUUID(id string) {
// 	if f.Struct.Name != "" {
// 		f.Struct.SetUUID(id)
// 	}
// }

// func (f Field) ForeignType(target string) any {
// 	if ft, ok := f.ForeignTypes[target]; ok {
// 		return ft
// 	}
// 	if ft, ok := ForeignTypes[target]; ok {
// 		return ft
// 	}
// 	return nil
// }

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

// Struct returns a Struct and error, but does not guarantee it is useful
func (f Field) Struct() (Struct, error) {
	return ToStruct(f.Value)
	// switch s, err := ToStruct(f.Value) ; {
	// 	case len(s.Fields())
	// }

}

// func (f Field) IsStruct() bool {
// 	return false
// }

// MultiValued returns true for any 'collection' type or any struct with more than one field
func (f Field) MultiValued() bool {
	kind := f.Value.Kind()

	switch {
	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind):
		return true
	case kind == reflect.Struct && len(f.Value.Children()) > 1:
		return true
	default:
		return false

	}
}
