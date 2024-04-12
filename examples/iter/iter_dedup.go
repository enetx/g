package main

import "github.com/enetx/g"

func main() {
	// Create a slice of integers with repeated elements
	g.SliceOf(1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		Dedup().   // Remove duplicate elements
		Collect(). // Collect the resulting slice
		Print()    // Print the collected slice: Slice[1, 3, 4, 8, 9]

	g.SliceOf([]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2}).
		Iter().
		Dedup().
		Collect().
		Print() // Slice[[1 2 3], [1 2]]
}
