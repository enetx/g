package cmp

import "cmp"

// Ordered represents an ordered value.
type Ordered int

// Then returns the receiver if it's non-zero, otherwise returns the other value.
func (o Ordered) Then(other Ordered) Ordered {
	if o != 0 {
		return o
	}

	return Ordered(other)
}

// Cmp compares two ordered values and returns the result as an Ordered value.
func Cmp[T cmp.Ordered](x, y T) Ordered { return Ordered(cmp.Compare(x, y)) }
