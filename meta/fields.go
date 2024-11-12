package meta

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

type Fields []Field

func ToFields(value any) (Fields, error) {
	fields, ok := value.(Fields)
	if ok {
		return fields, nil
	}
	s, err := ToStruct(value)
	if err != nil {
		return fields, err
	}
	return s.Fields(), nil

}

// Tags returns a map of Tags with keys that match field names
func (f Fields) Tags() map[string]Tags {
	fieldsTagMap := map[string]Tags{}
	for _, field := range f {
		fieldsTagMap[field.Name] = field.Tags()
	}
	return fieldsTagMap
}

// WithTag returns a subset of Fields with any of the given keys
func (f Fields) WithTag(keys ...any) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTag(keys...) {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithTagValue returns a subset of Fields with both the key and value
func (f Fields) WithTagValue(key, value string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTagValue(key, value) {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithTagTrue returns a subset of Fields whose Tags satisfy Tags.True
func (f Fields) WithTagTrue(keys ...any) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTagTrue(keys...) {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithTagFalse returns a subset of Fields whose Tags satisfy Tags.False
func (f Fields) WithTagFalse(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTagFalse(key) {
			// if field.Tags().False(key) {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithoutTag returns a subset of Fields that do not have any of the given keys
func (f Fields) WithoutTag(keys ...any) Fields {
	fields := Fields{}
	for _, field := range f {
		if !field.HasTag(keys...) {
			fields = append(fields, field)
		}
	}

	return fields
}

// WithoutTagValue returns a subset of Fields that do not have both the key and value
func (f Fields) WithoutTagValue(key, value string) Fields {
	fields := Fields{}
	for _, field := range f {
		if !field.HasTagValue(key, value) {
			fields = append(fields, field)
		}
	}
	return fields
}

// Names returns a slice of field names
func (f Fields) Names() []string {
	names := []string{}
	for _, field := range f {
		names = append(names, field.Name)
	}
	return names
}

// TagNames returns a slice of the field names according to field.TagName
func (f Fields) TagNames(keys ...any) []string {
	names := []string{}
	for _, field := range f {
		names = append(names, field.TagName(keys...))
	}
	return names
}

// Identifiers returns a slice of the fully namespaced identifiers
func (f Fields) Identifiers() []string {
	identifiers := []string{}
	for _, field := range f {
		identifiers = append(identifiers, field.Identifier())
	}
	return identifiers
}

// TagIdentifiers returns a slice of the fully namespaced identifiers
// from TagIdentifier
func (f Fields) TagIdentifiers(keys ...any) []string {
	identifiers := []string{}
	for _, field := range f {
		identifiers = append(identifiers, field.TagIdentifier(keys...))
	}
	return identifiers
}

// Types returns a slice of field types
func (f Fields) Types() []reflect.Type {
	types := []reflect.Type{}
	for _, field := range f {
		types = append(types, field.Type())
	}
	return types
}

func (f Fields) Kinds() []reflect.Kind {
	var kinds []reflect.Kind
	for _, field := range f {
		kinds = append(kinds, field.Kind())
	}
	return kinds
}

func (f Fields) MultiValued() Fields {
	var fields Fields
	for _, field := range f {
		if field.MultiValued() {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f Fields) ByNames(names ...any) Fields {
	nms := ToStringSlice(names...)
	var fields Fields
	for _, field := range f {
		if slices.Contains(nms, field.Name) {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f Fields) ByName(name string) Field {
	for _, field := range f {
		if field.Name == name {
			return field
		}
	}
	panic(fmt.Sprintf("Fields.ByName: name %s not found", name))
}

func (f Fields) ByTypes(types ...reflect.Type) Fields {
	var fields Fields
	for _, field := range f {
		if slices.Contains(types, field.Type()) {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f Fields) ByKinds(kinds ...reflect.Kind) Fields {
	var fields Fields
	for _, field := range f {
		if slices.Contains(kinds, field.Kind()) {
			fields = append(fields, field)
		}
	}
	return fields
}

// NonEmptyTagValues returns a slice of the first non-empty, non-nil tag that satisfies Tag.True and matches
// a key in the order given. In order to avoid slice length mismatches, it actually does use an empty string for fields
// where no match is found.
func (f Fields) NonEmptyTagValues(keys ...any) []string {
	var values []string
	for _, field := range f {
		values = append(values, field.NonEmptyTagValue(keys...))
	}
	return values
}

// // FirstTagValues returns a slice of the first Tag value for the given key
// func (f Fields) FirstTagValues(key string) []string {
// 	var values []string
// 	for _, field := range f {
// 		if tag := field.Tags().Tag(key); tag != nil {
// 			values = append(values, tag[0])
// 		}
// 	}
// 	return values
// }

// ElementsToStructs iterates through each field, attempts to convert it
// to a Struct, and returns a slice of the valid Structs. Note that it will
// panic if any field cannot be converted to a struct
func (f Fields) ToStructs() Structs {
	var s Structs
	for _, field := range f {
		str, err := field.ToStruct()
		if err != nil {
			continue
		}
		s = append(s, str)
	}
	return s
}

// Field is a convenience function to avoid having to index in a pipeline.
// It panics if Fields isn't length 1
func (f Fields) Field() Field {
	if len(f) == 1 {
		return f[0]
	}
	panic("Fields must be length 1")
}

func (f Fields) ToData() []Data {
	data := []Data{}
	for _, field := range f {
		data = append(data, field.ToData())
	}
	return data
}

func (f Fields) ToDataMap() map[string]Data {
	data := map[string]Data{}
	for _, field := range f {
		data[field.Name] = field.ToData()
	}
	return data
}

// TagTypes returns the 'type', prefering a tagged value, then falling back to reflect.Type
func (f Fields) TagTypes(keys ...any) []string {
	tagtypes := f.NonEmptyTagValues(keys...)
	ftypes := ToStringSlice(f.Types())
	return Coalesce(tagtypes, ftypes, "")
}

func (f Fields) TypeNames() []string {
	out := []string{}
	for _, field := range f {
		out = append(out, field.Type().Name())
	}
	return out
}

func (f Fields) TagTypeNames(keys ...any) []string {
	tagtypes := f.NonEmptyTagValues(keys...)
	ftypes := f.TypeNames()
	return Coalesce(tagtypes, ftypes, "")
}
