package filters

import "reflect"

// iszero represents an interface that defines the Empty() method.
type iszero interface {
	Empty() bool
}

// IsZero is a generic function designed to check if a value is considered zero.
// This function takes a value of type T and checks if it is "zero" according to the iszero interface
// or through a deep reflective comparison with the zero value of the type.
// If the value of type T implements the iszero interface, the function calls the Empty() method
// to determine if the value is "zero".
// If the type does not implement iszero, a deep reflective comparison is used with the zero value of type T.
func IsZero[T any](v T) bool {
	if value, ok := any(v).(iszero); ok {
		return value.Empty()
	}

	return reflect.DeepEqual(v, *new(T))
}
