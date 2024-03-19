package main

import (
	"github.com/enetx/g"
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
}
