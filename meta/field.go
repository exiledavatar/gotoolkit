package meta

import (
	"reflect"

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
	switch kind := ft.Kind(); kind {
	case reflect.Slice:
		return ft.Elem()
	case reflect.Array:
		return ft.Elem()
	case reflect.Map:
		return ft.Elem()
	default:
		return nil
	}
}

// Converts the element of an Array, Chan, Map, Pointer, or Slice to
// a Struct. Because it is designed for templating and piping, it panics
// if it cannot convert the field to a struct
func (f Field) ToStruct() (Struct, error) {
	var value any
	switch kind := f.Value.Kind(); {
	case (kind == reflect.Slice || kind == reflect.Array) && f.Value.Len() > 0:
		value = f.Value.Index(0).Interface()
	case kind == reflect.Map && f.Value.Len() > 0:
		value = f.Value.MapRange().Value()
	case (kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map) && f.Value.Len() == 0:
		value = reflect.New(f.Type().Elem()).Elem().Interface()
	default:
		value = f.Value.Interface()
	}
	s, err := NewStruct(value, Structconfig{
		Name:       f.Name,
		NameSpace:  f.Parent.NameSpace,
		Attributes: f.Attributes,
		Parent:     f.Parent,
		Tags:       f.Tags(),
	})
	if err != nil {
		return Struct{}, err
	}
	return s, nil
}

// Tags returns the parsed struct field tags
func (f Field) Tags() Tags {
	return ToTags(string(f.StructField.Tag))
}

// HasTag returns true any of the given tag keys exist
func (f Field) HasTag(keys ...string) bool {
	return f.Tags().Exists(keys...)
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
	for _, key := range keys {
		tag := f.Tags().Value(key)
		if len(tag) > 0 && tag.True() && tag[0] != "" {
			return tag[0]
		}
	}
	return f.Name
}

// Identifier uses the parent's Struct.Identifier and appends the field's
// Name to it
func (f Field) Identifier() string {
	return f.Parent.Identifier() + f.Parent.NameSpaceSeparator + f.Name
}

// TagIdentifier uses the parent's Struct.Identifier and appends the field's
// Name to it
func (f Field) TagIdentifier(keys ...string) string {
	return f.Parent.Identifier() + f.Parent.NameSpaceSeparator + f.TagName(keys...)
}

// Tag returns the tag for the given key, according to Tags.Tag
func (f Field) Tag(key string) Tag {
	return f.Tags().Tag(key)
}

// TagValueAtIndex returns the first tag.True, non-blank value for keys in the given order.
// If no match is found, it returns the empty string.
func (f Field) TagValueAtIndex(index int, keys ...string) string {
	for _, key := range keys {
		tag := f.Tags().Value(key)
		if len(tag) > 0 && tag.True() && tag[0] != "" {
			return tag[0]
		}
	}
	return f.Name
}

func (f Field) NonEmptyTagValue(keys ...string) string {
	if tag := f.Tags().NonEmptyValue(keys...); len(tag) > 0 {
		return tag[0]
	}
	return ""
}

// // Struct returns a Struct and error, but does not guarantee it is useful
// func (f Field) Struct() (Struct, error) {
// 	return ToStruct(f.Value)
// }

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

func (f Field) ToData() Data {
	value := f.Value.Interface()
	return Data(ToSlice(value))
}

// func (f Field) ToStructWithData() StructWithData {
// 	return StructWithData{
// 		Struct: f.ToStruct(),
// 		Data:   f.ToData(),
// 	}
// }
