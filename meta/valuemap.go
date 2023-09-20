package meta

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
)

// ValueMap's keys should match a struct's members' names (or tags) and values
// are wrapped in any()
type ValueMap map[string]any

func ToValueMap(value any, tagKey string) ValueMap {
	var fields Fields
	switch v := value.(type) {
	case Fields:
		fields = v
	case Struct:
		fields = v.Fields()
	default:
		str, err := ToStruct(value)
		if err != nil {
			return nil
		}
		fields = str.Fields()
	}
	if len(fields) == 0 {
		return nil
	}

	// to reduce the amount of tags necessary, we look for:
	// - if any tag.False(): include all except those fields
	// - if any tag.True():  include only those fields
	// - else: include all exported fields
	switch {
	case len(fields.WithTagFalse(tagKey)) > 0:
		// fields = fields.WithTagTrue(tagKey)
		var ff Fields
		for _, field := range fields {
			if !field.HasTagFalse(tagKey) {
				ff = append(ff, field)
			}
		}
		fields = ff
	case len(fields.WithTagTrue(tagKey)) > 0:
		fields = fields.WithTagTrue(tagKey)
	default:
		// do nothing - will include all fields
	}

	vm := ValueMap{}
	for _, field := range fields {
		var v any
		if field.Value.Kind() != reflect.Invalid {
			v = field.Value.Interface()
		}
		vm[field.Name] = v
	}
	return vm
}

// NewValueMap is a configurable version of ToValueMap - it converts a struct to a ValueMap
// based on the given tagKey, excludeValue, and includeValue
func NewValueMap(value any, tagKey, excludeValue, includeValue string) ValueMap {
	var fields Fields
	switch v := value.(type) {
	case Fields:
		fields = v
	case Struct:
		fields = v.Fields()
	default:
		str, err := ToStruct(value)
		if err != nil {
			return nil
		}
		fields = str.Fields()
	}
	if len(fields) == 0 {
		return nil
	}

	// to reduce the amount of tags necessary, we look for:
	// - if any tag.Contains(tagKey, excludeValue): include all except those fields
	// - if any tag.Contains(tagKey, includeValue):  include only those fields
	// - else: include all exported fields
	switch {
	case len(fields.WithTagValue(tagKey, excludeValue)) > 0:
		fields = fields.WithTagValue(tagKey, excludeValue)
	case len(fields.WithTagValue(tagKey, includeValue)) > 0:
		fields = fields.WithTagValue(tagKey, includeValue)
	default:
		// do nothing - will include all fields
	}

	vm := ValueMap{}
	for _, field := range fields {
		var v any
		if field.Value.Kind() != reflect.Invalid {
			v = field.Value.Interface()
		}
		vm[field.Name] = v
	}
	return vm

}

func (v ValueMap) Hash() string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", sha1.Sum(bb.Bytes()))
}
