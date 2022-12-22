package null

import "fmt"

type Nullable[T any] struct {
	V     T
	Valid bool
}

func New[T any](v T) Nullable[T] {
	return Nullable[T]{
		V:     v,
		Valid: true,
	}
}

func (n *Nullable[T]) Set(v T) {
	var x Nullable[T]
	if i := (interface{}(v)); i != nil {
		x = New[T](v)
	}
	*n = x
}

func (n Nullable[T]) Get() T {
	if n.Valid {
		return n.V
	}
	return *new(T)
}

func DispatchTest[T any](v T) {
	switch (interface{}(v)).(type) {
	case int:
		fmt.Println("type int")
	default:
		fmt.Println("other")
	}
}
