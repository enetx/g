package main

import "github.com/enetx/g"

func main() {
	// Example 1: Map each string in the slice to its uppercase version and print the result
	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Map(g.String.Upper). // Map each string to its uppercase version
		Collect().
		Print() // Slice[, BBB, DDD, , AAA, CCC]

	// Example 2: Map each string in the slice, replacing empty strings with "abc", and print the result
	g.SliceOf[g.String]("", "bbb", "ddd", "", "aaa", "ccc").
		Iter().
		Map(func(s g.String) g.String {
			if s.Empty() {
				s = "abc"
			}
			return s
		}).
		Collect().
		Print() // Slice[abc, bbb, ddd, abc, aaa, ccc]
}
