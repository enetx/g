package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a slice of integers with repeated elements
	g.SliceOf(1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		Dedup().   // Remove duplicate elements
		Collect(). // Collect the resulting slice
		Print()    // Print the collected slice
}
