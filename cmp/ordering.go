package cmp

import "fmt"

// Ordering is the result of a comparison between two values.
type Ordering int

const (
	Less    Ordering = iota - 1 // Less represents an ordered value where a compared value is less than another.
	Equal                       // Equal represents an ordered value where a compared value is equal to another.
	Greater                     // Greater represents an ordered value where a compared value is greater than another.
)

// Then returns the receiver if it's equal to Equal, otherwise returns the receiver.
func (o Ordering) Then(other Ordering) Ordering {
	if Equal == o {
		return Ordering(other)
	}

	return o
}

// Reverse reverses the ordering.
func (o Ordering) Reverse() Ordering {
	switch o {
	case Less:
		return Greater
	case Greater:
		return Less
	default:
		return Equal
	}
}

// String returns the string representation of the Ordering value.
func (o Ordering) String() string {
	switch o {
	case Less:
		return "Less"
	case Equal:
		return "Equal"
	case Greater:
		return "Greater"
	default:
		return fmt.Sprintf("Unknown Ordering value: %d", int(o))
	}
}
