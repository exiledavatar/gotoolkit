package meta

import "strings"

type Namespace string

func (n Namespace) Separator() string {
	return "."
}

func (n Namespace) Join(ns ...string) Namespace {

	out := strings.Join(ns, ".")
	// x := append([]string{string(n)}, ns...)
	return Namespace(out)
}

// func (n Namespace) Test() {

// 	i := 0

// }
type Name struct {
	name      string
	namespace []string
	separator string
}

func ParseName(n string, sep string) Name {
	fn := strings.Split(n, sep)
	switch {
	case len(fn) == 1:
		return Name{
			name:      fn[0],
			separator: sep,
		}
	case len(fn) > 1:
		return Name{
			name:      fn[0],
			namespace: fn[1:],
			separator: sep,
		}
	default:
		return Name{
			separator: sep,
		}
	}
}

func ToName(n string) Name {
	return ParseName(n, ".")
}

func (n Name) String() string {
	fn := append([]string{n.name}, n.namespace...)
	return strings.Join(fn, n.separator)
}

func (n *Name) SetNamespace(ns ...string) Name {
	n.namespace = ns
	return *n
}

func (n *Name) SetSeparator(sep string) Name {
	n.separator = sep
	return *n
}

