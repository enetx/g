package cmp

import "cmp"

// Cmp compares two ordered values and returns the result as an Ordering value.
func Cmp[T cmp.Ordered](x, y T) Ordering { return Ordering(cmp.Compare(x, y)) }
