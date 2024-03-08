package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	// Create a slice of strings with repeated elements
	chunks := g.SliceOf[g.String]("bbb", "ddd", "bbb", "ccc", "aaa", "bbb", "ccc").
		Iter().
		Unique().  // Remove duplicates from the slice
		Chunks(2). // Split the slice into chunks of size 2
		Collect()

	// Print the resulting chunks
	fmt.Println(chunks)
}
