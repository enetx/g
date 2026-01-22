package main

import . "github.com/enetx/g"

func main() {
	// Example 1: Map each string in the slice to its uppercase version and print the result
	SliceOf[String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Map(String.Upper). // Map each string to its uppercase version
		Collect().
		Println() // Slice[, BBB, DDD, , AAA, CCC]

	// Example 2: Map each string in the slice, replacing empty strings with "abc", and print the result
	SliceOf[String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Map(func(s String) String {
			if s.IsEmpty() {
				s = "abc"
			}
			return s
		}).
		Collect().
		Println() // Slice[abc, bbb, ddd, abc, aaa, ccc]
}
