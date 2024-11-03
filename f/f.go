package f

import (
	"cmp"
	"reflect"

	"github.com/enetx/g/pkg/constraints"
)

// Comparable reports whether the value v is comparable.
func Comparable[T any](v T) bool { return reflect.ValueOf(v).Comparable() }

// Zero is a generic function designed to check if a value is considered zero.
func Zero[T cmp.Ordered](v T) bool { return v == *new(T) }

// Even is a generic function that checks if the provided integer is even.
func Even[T constraints.Integer](i T) bool { return i%2 == 0 }

// Odd is a generic function that checks if the provided integer is odd.
func Odd[T constraints.Integer](i T) bool { return i%2 != 0 }

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

// Eqd returns a comparison function that evaluates to true when a value is deeply equal to the provided threshold.
func Eqd[T any](t T) func(T) bool {
	return func(s T) bool {
		return reflect.DeepEqual(t, s)
	}
}

// Ned returns a comparison function that evaluates to true when a value is not deeply equal to the provided threshold.
func Ned[T any](t T) func(T) bool {
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

// Gte returns a comparison function that evaluates to true when a value is greater than or equal to the threshold.
func Gte[T cmp.Ordered](t T) func(T) bool {
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

// Lte returns a comparison function that evaluates to true when a value is less than or equal to the threshold.
func Lte[T cmp.Ordered](t T) func(T) bool {
	return func(s T) bool {
		return s <= t
	}
}
