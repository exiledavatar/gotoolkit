package meta

// Config allows default config options for various functions
// Note: this is not yet implemented
type Config struct {
	NameSpace                []string
	NameSpaceSeparator       string // default is "."
	UUID                     func(any) string
	Tags                     Tags           // struct tags, not field tags, may optionally be parsed from tags labeled struct
	RemoveExistingTags       bool           // remove existing tags - false will simply apopend any included tags to current ones
	Attributes               map[string]any // the kitchen sink... for everything else
	RemoveExistingAttributes bool           // remove existing constraints - false will simply apopend any included constraints to current ones
}

var Defaults = Config{}
