package meta

import (
	"regexp"
	"strings"
)

// Tags are the parsed struct field tags as
// map[tagLabel][]tagvalues
type Tags map[string]Tag

func ToTags(s string) Tags {
	pattern := regexp.MustCompile(`(?m)(?P<key>\w*):\"(?P<value>[^"]*)\"`)
	matches := pattern.FindAllStringSubmatch(s, -1)
	var tkv = map[string]Tag{}
	for _, match := range matches {
		tkv[match[1]] = strings.Split(match[2], ",")
	}
	return tkv
}

// Value returns the parsed Tag for the first key found, in order given, or nil if missing
// note: it is currently the same as Tags.Tag but may be changed to return the actual string value
func (t Tags) Value(keys ...any) Tag {
	ks := ToStringSlice(keys...)
	for _, key := range ks {
		if tag, ok := t[key]; ok && tag.True() {
			return tag
		}
	}
	return nil
}

// NonEmptyValue returns the parsed Tag for the first key found, in order given,
// where the tag is non-nil, non-empty, and satisfies tag.True.
// It returns nil if no match is found.
func (t Tags) NonEmptyValue(keys ...any) Tag {
	ks := ToStringSlice(keys...)
	for _, key := range ks {
		if tag, ok := t[key]; ok && len(tag) > 0 && tag.True() && tag[0] != "" {
			return tag
		}
	}
	return nil
}

// ByKeys returns a subset of Tags with the given keys, or nil if none are found
func (t Tags) ByKeys(keys ...any) Tags {
	ks := ToStringSlice(keys...)
	tags := Tags{}
	for _, key := range ks {
		if value, ok := t[key]; ok {
			tags[key] = value
		}
	}
	if len(tags) == 0 {
		return nil
	}
	return tags
}

// False only returns true if the tags exists and the first value matches
// one in ConfigTagFalse (by default this is just "-")
func (t Tags) False(keys ...any) bool {
	ks := ToStringSlice(keys...)
	var out bool
	for _, key := range ks {
		switch tag, ok := t[key]; {
		case out:
			return out
		case ok && tag != nil:
			out = tag.False()
		}
	}
	return false
}

// True returns true if the tag exists and the first value does not match
// one in ConfigTagFalse (by default this is just "-")
func (t Tags) True(keys ...any) bool {
	ks := ToStringSlice(keys...)
	var out bool
	for _, key := range ks {
		switch tag, ok := t[key]; {
		case out:
			return out
		case ok && tag != nil:
			out = tag.True()
		}
	}
	return false
}

// Exists returns true if a tag with any of the given keys exists,
// even if it is empty
func (t Tags) Exists(keys ...any) bool {
	ks := ToStringSlice(keys...)
	for _, key := range ks {
		if tag, ok := t[key]; ok && tag != nil {
			return true
		}
	}
	return false
}

// Contains returns true if the tag with key has the value,
// regardless of its index in Tag
func (t Tags) Contains(key, value string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.Contains(value)
	}
	return false
}

// NotContains returns true if there is no tag with the key and value
func (t Tags) NotContains(key, value string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.NotContains(value)
	}
	return true
}

// Tag returns the tag for key, or nil if it is missing
func (t Tags) Tag(keys ...any) Tag {
	ks := ToStringSlice(keys...)
	for _, key := range ks {
		if tag, ok := t[key]; ok && tag.True() {
			return tag
		}
	}
	return nil
}

// Set is a convient method wrapper for assigning (replacing) a key's tags
func (t Tags) Set(key string, values ...string) Tags {
	t[key] = values
	return t
}

// Prepend inserts the value to the beginning of the key's tags
// it is safe to use even if the key doesn't currently exist
func (t Tags) Prepend(key, value string) Tags {
	tag := []string{value}
	tag = append(tag, t[key]...)
	t[key] = tag
	return t
}

// Append adds the value to the end of the key's tags
// it is safe to use even if the key doesn't currently exist
func (t Tags) Append(key, value string) Tags {
	switch tag, ok := t[key]; {
	case ok:
		t[key] = append(tag, value)
	default:
		t[key] = []string{value}
	}
	return t
}

func (t Tags) Replace(key, value string, index int) Tags {

	switch tag, ok := t[key]; {
	case ok && len(tag) >= index+1:
		tag[index] = value
		t[key] = tag
	case ok:

		length := len(tag)
		if length <= index {
			length = index + 1
		}
		dst := make([]string, length)
		copy(dst, tag)
		dst[index] = value
		t[key] = dst
	default:
		tag = make([]string, index+1)
		tag[index] = value
		t[key] = tag
	}
	return t
}

func (t Tags) Insert(key, value string, index int) Tags {
	switch tag, ok := t[key]; {
	case ok:
		dst := []string{}
		for i := 0; i <= len(tag) || i <= index; i++ {
			switch {
			case i == index && i < len(tag):
				dst = append(dst, value, tag[i])
			case i == index:
				dst = append(dst, value)
			case i < len(tag):
				dst = append(dst, tag[i])
			default:
				dst = append(dst, "")
			}
		}
		t[key] = dst
	default:
		tag = make([]string, index+1)
		tag[index] = value
		t[key] = tag
	}
	return t
}
