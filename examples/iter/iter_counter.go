package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	slice := SliceOf(1, 2, 3, 1, 2, 1)

	slice.Iter().
		Counter().
		SortBy(func(a, b Pair[int, Int]) cmp.Ordering { return b.Value.Cmp(a.Value) }).
		Collect().
		Println() // MapOrd{1:3, 2:2, 3:1}
}
