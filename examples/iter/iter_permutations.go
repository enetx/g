package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a slice of integers
	slice := Slice[int]{1, 2, 3}

	// Get an iterator for the slice, generate permutations, and collect the result
	perms := slice.Iter().Permutations().Collect()

	// Iterate over the permutations and print each one
	for _, perm := range perms {
		fmt.Println(perm)
	}

	// Create two slices of strings
	p1 := SliceOf[String]("bbb", "ddd")
	p2 := SliceOf[String]("xxx", "aaa")

	// Chain the two slices, convert to uppercase, generate permutations, and collect the result
	pp := p1.
		Iter().
		Chain(p2.Iter()).
		Map(String.Upper).
		Permutations().
		Collect()

	// Iterate over the permutations and print each one
	for _, v := range pp {
		v.Println()
	}
}
