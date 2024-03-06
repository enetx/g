package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a slice of strings
	ps := g.SliceOf[g.String]("bbb", "ddd", "xxx", "aaa", "ccc").
		Iter().
		Enumerate(). // Enumerate the elements, creating pairs of index and value
		Collect()    // Collect the resulting pairs

	// Print the ordered map
	ps.Print()
}
