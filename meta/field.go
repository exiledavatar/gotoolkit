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
	Attributes map[string]any //string // catchall for additional attributes
	Parent     *Struct
	Value
	reflect.StructField
}

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

// Converts the element of a Array, Chan, Map, Pointer, or Slice to
// a Struct.
func (f Field) ElementToStruct() (Struct, error) {

	s, err := NewStruct(f.Value.Type().Elem(), Structconfig{
		Name:       f.Name,
		NameSpace:  f.Parent.NameSpace,
		Attributes: f.Attributes,
		Tags:       f.Tags(), //any(f.Tags()).(map[string][]string),
	})
	if err != nil {
		return s, err
	}
	return s, nil
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

// TagName ranges through the provided keys in order and returns the first non-blank, non-false value, or field.Name if none are found.
func (f Field) TagName(keys ...string) string {
	name := f.Name
	for _, key := range keys {
		switch tag := f.Tags().Tag(key); {
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

// Identifier uses the parent's Struct.Identifier and appends the field's
// Name to it
func (f Field) Identifier() string {
	name := strings.ToLower(f.Name)
	return f.Parent.Identifier() + f.Parent.NameSpaceSeparator + name
}

// TagIdentifier uses the parent's Struct.Identifier and appends the field's
// Name to it
func (f Field) TagIdentifier(keys ...string) string {
	name := strings.ToLower(f.TagName(keys...))
	return f.Parent.Identifier() + f.Parent.NameSpaceSeparator + name
}

func (f Field) Tag(key string) Tag {
	return f.Tags().Tag(key)
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

// Struct returns a Struct and error, but does not guarantee it is useful
func (f Field) Struct() (Struct, error) {
	return ToStruct(f.Value)
}

// MultiValued returns true for any 'collection' type or any struct with more than one field
func (f Field) MultiValued() bool {
	switch kind := f.Value.Kind(); {
	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind):
		return true
	case kind == reflect.Struct && len(f.Value.Children()) > 1:
		return true
	default:
		return false

	}
}
