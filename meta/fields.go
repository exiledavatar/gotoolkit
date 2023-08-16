package meta

import (
	"log"
	"reflect"

	"golang.org/x/exp/slices"
)

type Fields []Field

func ToFields(value any) (Fields, error) {
	// var fields Fields
	fields, ok := value.(Fields)
	if ok {
		log.Println("already Fields")
		return fields, nil
	}

	s, err := ToStruct(value)
	if err != nil {
		log.Println("ToFields - ToStruct err not nil")
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

// WithTag returns a subset of Fields with the key
func (f Fields) WithTag(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTag(key) {
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
func (f Fields) WithTagTrue(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTagTrue(key) {
			// if field.Tags().True(key) {
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

// WithoutTag returns a subset of Fields that do not have the key
func (f Fields) WithoutTag(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if !field.HasTag(key) {
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
func (f Fields) TagNames(keys ...string) []string {
	names := []string{}
	for _, field := range f {
		names = append(names, field.TagName(keys...))
	}
	return names
}

// Identifiers returns a slice of the fully namespaced identifiers
func (f Fields) Identifiers(key string) []string {
	identifiers := []string{}
	for _, field := range f {
		identifiers = append(identifiers, field.Identifier())
	}
	return identifiers
}

// TagIdentifiers returns a slice of the fully namespaced identifiers
// from TagIdentifier
func (f Fields) TagIdentifiers(key string) []string {
	identifiers := []string{}
	for _, field := range f {
		identifiers = append(identifiers, field.TagIdentifier(key))
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

func (f Fields) ByNames(names ...string) Fields {
	var fields Fields
	for _, field := range f {
		if slices.Contains(names, field.Name) {
			fields = append(fields, field)
		}
	}
	return fields

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

// FirstTagValues returns a slice of the first Tag value for the given key
func (f Fields) FirstTagValues(key string) []string {
	var values []string
	for _, field := range f {
		if tag := field.Tags().Tag(key); tag != nil {
			values = append(values, tag[0])
		}
	}
	return values
}

// ElementsToStructs iterates through each field, attempts to convert it
// to a Struct, and returns a slice of the valid Structs.
func (f Fields) ElementsToStructs() Structs {
	var s Structs
	for _, field := range f {
		if st, err := field.ElementToStruct(); err == nil {
			s = append(s, st)
		}
	}
	return s
}
