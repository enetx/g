package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	// Create a slice of integers
	windows := g.SliceOf(1, 2, 3, 4).
		Iter().
		Windows(2). // Create windows of size 2
		Collect()   // Collect the resulting windows

	// Print the collected windows
	fmt.Println(windows)

	// Convert to iterator
	g.SliceOf(windows...).Iter()
}
