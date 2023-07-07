package meta

import (
	"reflect"
)

type Fields []Field

func ToFields(value any) (Fields, error) {
	// var fields Fields
	fields, ok := value.(Fields)
	if ok {
		return fields, nil
	}

	s, err := ToStruct(value)
	if err != nil {
		return fields, err
	}

	var sfs []reflect.StructField
	for _, sf := range reflect.VisibleFields(s.Type()) {
		if sf.Anonymous || !sf.IsExported() {
			continue
		}
		sfs = append(sfs, sf)
	}

	for i, child := range s.Children() {
		field := Field{
			Name:        child.Name,
			StructField: sfs[i],
			Value:       child,
			Parent:      &s,
		}
		fields = append(fields, field)

	}

	return fields, nil
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

func (f Fields) MultiValued() Fields {
	var fields Fields
	for _, field := range f {
		if field.Value.Kind() == reflect.Struct && field.MultiValued() {
			fields = append(fields, field)
		}
	}
	return fields
}

// func (f Fields) SetUUID(id string) {
// 	for i, _ := range f {
// 		f[i].SetUUID(id)
// 	}
// }

// func (f Fields) SetUUIDPrefix(id string) {
// 	for i, _ := range f {
// 		f[i].SetUUID(id)
// 	}
// }

// func (f Fields) SetUUIDSuffix(suffix string) {
// 	for i, _ := range f {
// 		id := f[i].UUID
// 		f[i].SetUUID()
// 	}
// }
