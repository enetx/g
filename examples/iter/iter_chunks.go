package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a slice of strings with repeated elements
	chunks := SliceOf[String]("bbb", "ddd", "bbb", "ccc", "aaa", "bbb", "ccc").
		Iter().
		Unique().  // Remove duplicates from the slice
		Chunks(2). // Split the slice into chunks of size 2
		Collect()

	// Print the resulting chunks
	fmt.Println(chunks)
}
