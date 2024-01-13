package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a slice of strings
	g.SliceOf[g.String]("bbb", "ddd", "bbb", "aaa", "bbb").
		Iter().
		Unique().  // Filter out duplicate elements
		Collect(). // Collect the unique elements
		Print()    // Print the collected unique elements
}
