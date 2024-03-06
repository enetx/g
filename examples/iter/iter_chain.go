package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	// Create two slices of strings, p1 and p2
	p1 := g.SliceOf[g.String]("bbb", "ddd")
	p2 := g.SliceOf[g.String]("xxx", "aaa")

	// Chain the iterators of p1 and p2 and collect the results into a new slice pp
	pp := p1.
		Iter().
		Chain(p2.Iter()).
		Collect()

	// Iterate over the resulting slice pp and print each element
	for _, v := range pp {
		v.Print()
	}

	// iter1 := g.NewMap[int, string]().Set(1, "a").Iter()
	// iter2 := g.NewMap[int, string]().Set(2, "b").Iter()

	// iter1.Chain(iter2).Collect().Print()

	// set1 := g.NewSet[int]().Add(1, 2, 3).Iter()
	// set2 := g.NewSet[int]().Add(3, 3, 2, 3, 4, 5).Iter()

	// set1.Chain(set2).Collect().Print()
}