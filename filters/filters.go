package filters

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

// IsEq is a generic function that checks if two comparable values are equal using the equality operator.
func IsEq[T comparable](x, y T) bool { return x == y }

// IsEqDeep is a generic function that checks if two any values are deeply equal using reflect.DeepEqual.
func IsEqDeep[T any](x, y T) bool { return reflect.DeepEqual(x, y) }
