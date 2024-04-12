package main

import "github.com/enetx/g"

func main() {
	slice := g.SliceOf(1, 2, 3, 1, 2, 1)

	slice.Iter().
		Counter().
		SortBy(func(a, b g.Pair[int, uint]) bool { return a.Value > b.Value }).
		Collect().
		Print() // MapOrd{1:3, 2:2, 3:1}
}
