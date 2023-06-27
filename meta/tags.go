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

// Value returns the parsed Tag for the given key, or nil if missing
func (t Tags) Value(key string) Tag {
	if value, ok := t[key]; ok {
		return value
	}
	return nil
}

// False only returns true if the tags exists and the first value matches
// one in ConfigTagFalse (by default this is just "-")
func (t Tags) False(key string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.False()
	}
	return false
}

// True returns true if the tag exists and the first value does not match
// one in ConfigTagFalse (by default this is just "-")
func (t Tags) True(key string) bool {
	if tag, ok := t[key]; ok && tag != nil {
		return tag.True()
	}
	return false
}

// Exists returns true if the tag with key exists, even if it is empty
func (t Tags) Exists(key string) bool {
	tag, ok := t[key]
	return ok && tag != nil
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
// note: it is the same as Value
func (t Tags) Tag(key string) Tag {
	if tag, ok := t[key]; ok {
		return tag
	}
	return nil
}
