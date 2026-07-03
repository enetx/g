package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a slice of integers
	slice := Slice[int]{0, 1, 2, 3}

	// Get an iterator for the slice, generate combinations, and collect the result
	perms := slice.Iter().Combinations(3).Collect()

	// Iterate over the combinations and print each one
	// 012 013 023 123
	for _, perm := range perms {
		fmt.Println(perm)
	}

	// Create two slices of strings
	p1 := SliceOf[String]("a", "b")
	p2 := SliceOf[String]("c", "d")

	// Chain the two slices, convert to uppercase, generate combinations, and collect the result
	pp := p1.
		Iter().
		Chain(p2.Iter()).
		Map(String.Upper).
		Combinations(2).
		Collect()

	// Iterate over the combinations and print each one
	// AB AC AD BC BD CD
	for _, v := range pp {
		v.Println()
	}
}
