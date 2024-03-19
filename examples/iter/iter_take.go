package main

import "github.com/enetx/g"

func main() {
	// Create a slice of strings with empty and non-empty values
	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Exclude(g.String.Empty).
		Cycle().   // Cycle through the elements in the slice
		Take(20).  // Take the first 20 elements from the cycled sequence
		Collect(). // Collect the resulting sequence
		Print()    // Print the collected sequence
}
