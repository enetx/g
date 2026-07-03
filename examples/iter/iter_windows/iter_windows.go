package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a slice of integers
	windows := SliceOf(1, 2, 3, 4).
		Iter().
		Windows(2). // Create windows of size 2
		Collect()   // Collect the resulting windows

	// Print the collected windows
	fmt.Println(windows) // [Slice[1, 2] Slice[2, 3] Slice[3, 4]]

	// Convert to iterator
	SliceOf(windows...).Iter()
}
