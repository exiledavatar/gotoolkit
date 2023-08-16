package meta

import (
	"golang.org/x/exp/slices"
)

type Tag []string

var ConfigTagFalse = []string{"-"}

// Contains is a wrapper for slices.Contains that returns true if
// the value is explicitly in the tag
func (t Tag) Contains(values ...string) bool {
	for _, value := range values {
		if slices.Contains(t, value) {
			return true
		}
	}
	return false
}

// NotContains is a wrapper for slices.Contains that returns true if
// the value isn't in the tag
func (t Tag) NotContains(values ...string) bool {
	for _, value := range values {
		if slices.Contains(t, value) {
			return false
		}
	}
	return true
}

// Index is a wrapper for slices.Index that returns the index
// of value in the tag, or -1 if not present
func (t Tag) Index(values ...string) int {
	var index int
	for _, value := range values {
		index = slices.Index(t, value)
		if index != -1 {
			return index
		}
	}
	return index
}

// False only returns true if the tag's first value matches
// one in ConfigTagFalse (by default this is just "-")
func (t Tag) False() bool {
	return slices.Contains(ConfigTagFalse, t[0])
}

// True returns true if the first value doesn't match
// one in ConfigTagFalse (by default this is just "-")
func (t Tag) True() bool {
	return !slices.Contains(ConfigTagFalse, t[0])
}
