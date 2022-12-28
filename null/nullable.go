package null

import (
	"encoding/json"
	"fmt"
)

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

func FromPointer[T any](pointer *T) Nullable[T] {
	return New(*pointer)
}

func (n *Nullable[T]) Set(v T) {
	var x Nullable[T]
	if i := (interface{}(v)); i != nil {
		x = New(v)
	}
	*n = x
}

func (n Nullable[T]) Get() T {
	if n.Valid {
		return n.V
	}
	return *new(T)
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.V)
}

// TODO: write UnmarshalJSON
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	var v T
	var s string
	switch err := json.Unmarshal(b, &s); {
	case b == nil, err != nil, s == "null":
		*n = Nullable[T]{V: v, Valid: false}
		return nil
	default:
		err := json.Unmarshal(b, &v)
		*n = New(v)
		return err
	}
}

// TODO: is there a general purpose scan function like json.Unmarshal?
// func (n *Nullable[T]) Scan(src any) error {
// 	if src == nil {
// 		n.V = *new(T)
// 		n.Valid = false
// 		return nil
// 	}
// 	*n = New()

// 	return nil
// }

// TODO: write Value
// func (n Nullable[T]) Value() (driver.Value, error) {
// 	return nil, nil
// }

func (n Nullable[T]) MarshalText() ([]byte, error) {
	// if !n.Valid {
	// 	return []byte{}, nil
	// }
	// if i, ok := (interface{}(n.V)).(encoding.TextMarshaler); ok {
	// 	return i.MarshalText()
	// }
	// return []byte(fmt.Sprint(n.V)), nil
	return n.MarshalJSON()
}

// TODO: is there a general purpose text unmarshaling function like json.Unmarshal
func (n *Nullable[T]) UnmarshalText(b []byte) error {
	return n.UnmarshalJSON(b)
}

// should probably leave gostring alone
// func (n Nullable[T]) GoString() string {
// 	if !n.Valid {
// 		return "null"
// 	}
// 	return fmt.Sprintf("%#v", n.Value)
// }

func (n Nullable[T]) String() string {
	if !n.Valid {
		return "null"
	}
	return fmt.Sprint(n.V)
}

// func (n *Nullable[T]) ()
// func (n *Nullable[T]) ()

func DispatchTest[T any](v T) {
	switch (interface{}(v)).(type) {
	case int:
		fmt.Println("type int")
	default:
		fmt.Println("other")
	}
}
