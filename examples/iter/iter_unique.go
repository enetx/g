package main

import "github.com/enetx/g"

func main() {
	// Create a slice of strings
	g.SliceOf[g.String]("bbb", "ddd", "bbb", "aaa", "xxx", "bbb", "bbb", "xxx", "ddd", "bbb", "aaa", "bbb").
		Iter().
		Unique().  // Filter out duplicate elements
		Collect(). // Collect the unique elements
		Print()    // Print the collected unique elements
}
