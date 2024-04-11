package f

import (
	"cmp"
	"reflect"

	"github.com/enetx/g/pkg/constraints"
)

// IsZero is a generic function designed to check if a value is considered zero.
func IsZero[T cmp.Ordered](v T) bool { return v == *new(T) }

// IsEven is a generic function that checks if the provided integer is even.
func IsEven[T constraints.Integer](int T) bool { return int%2 == 0 }

// IsOdd is a generic function that checks if the provided integer is odd.
func IsOdd[T constraints.Integer](int T) bool { return int%2 != 0 }

// Eq returns a comparison function that evaluates to true when a value is equal to the provided threshold.
func Eq[T comparable](t T) func(T) bool {
	return func(s T) bool {
		return s == t
	}
}

// Ne returns a comparison function that evaluates to true when a value is not equal to the provided threshold.
func Ne[T comparable](t T) func(T) bool {
	return func(s T) bool {
		return s != t
	}
}

// EqDeep returns a comparison function that evaluates to true when a value is deeply equal to the provided threshold.
func EqDeep[T any](t T) func(T) bool {
	return func(s T) bool {
		return reflect.DeepEqual(t, s)
	}
}

// NeDeep returns a comparison function that evaluates to true when a value is not deeply equal to the provided threshold.
func NeDeep[T any](t T) func(T) bool {
	return func(s T) bool {
		return !reflect.DeepEqual(t, s)
	}
}

// Gt returns a comparison function that evaluates to true when a value is greater than the threshold.
func Gt[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s > t
	}
}

// GtEq returns a comparison function that evaluates to true when a value is greater than or equal to the threshold.
func GtEq[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s >= t
	}
}

// Lt returns a comparison function that evaluates to true when a value is less than the threshold.
func Lt[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s < t
	}
}

// LtEq returns a comparison function that evaluates to true when a value is less than or equal to the threshold.
func LtEq[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s <= t
	}
}
