package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Create a slice of strings with empty and non-empty values
	SliceOf("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Exclude(f.IsZero).
		Cycle().   // Cycle through the elements in the slice
		Take(10).  // Take the first 10 elements from the cycled sequence
		Collect(). // Collect the resulting sequence
		Println()  // Print the collected sequence

	// Output: Slice[bbb, ddd, aaa, ccc, bbb, ddd, aaa, ccc, bbb, ddd]

	// Predicate-based slicing: TakeWhile stops at the first false,
	// SkipWhile drops the leading run and yields the rest
	SliceOf(1, 2, 3, 10, 4).
		Iter().
		TakeWhile(func(v int) bool { return v < 5 }).
		Collect().
		Println() // Slice[1, 2, 3]

	SliceOf(1, 2, 3, 10, 4).
		Iter().
		SkipWhile(func(v int) bool { return v < 5 }).
		Collect().
		Println() // Slice[10, 4]

}
