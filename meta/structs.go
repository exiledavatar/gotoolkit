package meta

type Structs []Struct

func (s Structs) TagName(keys ...string) []string {
	var names []string
	for _, str := range s {
		names = append(names, str.TagName(keys...))
	}
	return names
}

func (s Structs) Identifiers() []string {
	var identifiers []string
	for _, str := range s {
		identifiers = append(identifiers, str.Identifier())
	}
	return identifiers
}

func (s Structs) ToStructMap() map[string]Struct {
	structmap := map[string]Struct{}
	for _, str := range s {
		structmap[str.Name] = str
	}
	return structmap
}
