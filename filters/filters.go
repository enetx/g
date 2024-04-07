package filters

import (
	"cmp"

	"github.com/enetx/g/pkg/constraints"
)

// IsZero is a generic function designed to check if a value is considered zero.
func IsZero[T cmp.Ordered](v T) bool { return v == *new(T) }

// IsEven is a generic function that checks if the provided integer is even.
func IsEven[T constraints.Integer](int T) bool { return int%2 == 0 }

// IsOdd is a generic function that checks if the provided integer is odd.
func IsOdd[T constraints.Integer](int T) bool { return int%2 != 0 }
