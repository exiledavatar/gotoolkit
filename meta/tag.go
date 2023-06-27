package meta

import (
	"golang.org/x/exp/slices"
)

type Tag []string

var ConfigTagFalse = []string{"-"}

// Contains is a wrapper for slices.Contains that returns true if
// the value is explicitly in the tag
func (t Tag) Contains(value string) bool {
	return slices.Contains(t, value)
}

// NotContains is a wrapper for slices.Contains that returns true if
// the value isn't in the tag
func (t Tag) NotContains(value string) bool {
	return !slices.Contains(t, value)
}

// Index is a wrapper for slices.Index that returns the index
// of value in the tag, or -1 if not present
func (t Tag) Index(value string) int {
	return slices.Index(t, value)
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
