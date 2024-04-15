package cmp

import "cmp"

// Ordered represents an ordered value.
type Ordered int

const (
	Less    Ordered = iota - 1 // Less represents an ordered value where a compared value is less than another.
	Equal                      // Equal represents an ordered value where a compared value is equal to another.
	Greater                    // Greater represents an ordered value where a compared value is greater than another.
)

// Then returns the receiver if it's equal to Equal, otherwise returns the receiver.
func (o Ordered) Then(other Ordered) Ordered {
	if o == Equal {
		return Ordered(other)
	}
	return o
}

// Reverse returns the reverse of the ordered value.
func (o Ordered) Reverse() Ordered {
	switch o {
	case Less:
		return Greater
	case Greater:
		return Less
	default:
		return Equal
	}
}

// Cmp compares two ordered values and returns the result as an Ordered value.
func Cmp[T cmp.Ordered](x, y T) Ordered { return Ordered(cmp.Compare(x, y)) }
