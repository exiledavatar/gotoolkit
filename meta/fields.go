package meta

import (
	"reflect"
)

type Fields []Field

// ToFields returns the exported fields of a struct, pointer to a struct,
// reflect.Value of a struct, or reflect.Type of a struct
func ToFields(value any) Fields {
	var fields Fields

	rv, rt, _ := ToIndirectReflectValue(value)
	if !rv.IsValid() {
		return fields
	}

	if rt.Kind() != reflect.Struct {
		return fields
	}

	sfs := reflect.VisibleFields(rt)
	for _, sf := range sfs {
		if !sf.IsExported() || sf.Anonymous {
			continue
		}
		rfv, rft, rfPointer := ToIndirectReflectValue(rv.FieldByName(sf.Name))

		switch {
		case rfv.Kind() == reflect.Slice && rfv.Len() > 0:
			// take rfv.Index(0)....
		case rfv.Kind() == reflect.Slice && rfv.Len() < 1:
			// ?rft.Elem()
		case rfv.Kind() == reflect.Struct:
			// use ToStruct as normal?
		case rfv.Kind() == reflect.Map:
			// use
		}
		// if rfv.Kind() == reflect.Slice {
		// 	rfv.Elem()
		// }
		s, _ := ToStruct(rfv)

		sf.Type = rft
		field := Field{
			Name:        sf.Name,
			StructField: sf,
			Value:       rfv,
			pointer:     rfPointer,
			Struct:      s,
		}
		fields = append(fields, field)
	}
	return fields
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

// Types returns a slice of field types
func (f Fields) Types() []reflect.Type {
	types := []reflect.Type{}
	for _, field := range f {
		types = append(types, field.Type())
	}
	return types
}

func (f Fields) SetUUID(id string) {
	for i, _ := range f {
		f[i].SetUUID(id)
	}
}
