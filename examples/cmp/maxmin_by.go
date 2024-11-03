package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	mini := cmp.MinBy(cmp.Cmp, 3, 1, 4, 2, 5)
	fmt.Println("minimum int:", mini)

	mins := cmp.MinBy(cmp.Cmp, "banana", "apple", "orange")
	fmt.Println("minimum string:", mins)

	maxi := cmp.MaxBy(cmp.Cmp, 3, 1, 4, 2, 5)
	fmt.Println("maximum integer:", maxi)

	maxs := cmp.MaxBy(cmp.Cmp, "banana", "apple", "orange")
	fmt.Println("maximum string:", maxs)

	// cmp function
	ord := func(a, b []Int) cmp.Ordering { return a[0].Cmp(b[0]) }

	maxsis := cmp.MaxBy(ord, [][]Int{{1, 2, 3, 4}, {5, 6, 7, 8}}...)
	fmt.Printf("maximum []g.Int: %v\n", maxsis)

	maxgis := SliceOf([][]Int{{1, 2, 3, 4}, {5, 6, 7, 8}}...).MaxBy(ord)
	fmt.Printf("maximum []g.Int: %v\n", maxgis)

	maxgsi := SliceOf(1, 2, 3, 4, 5).MaxBy(cmp.Cmp)
	fmt.Printf("maximum int: %v\n", maxgsi)
}
