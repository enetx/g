package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Create two slices of strings, p1 and p2
	p1 := SliceOf[String]("bbb", "ddd")
	p2 := SliceOf[String]("xxx", "aaa")

	// Chain the iterators of p1 and p2 and collect the results into a new slice pp
	pp := p1.
		Iter().
		Parallel(10).
		Chain(p2.Iter().Parallel().Map(String.Upper).Filter(f.Eq(String("AAA")))).
		Collect()

	// Iterate over the resulting slice pp and print each element
	for _, v := range pp {
		v.Println()
	}
}
