package main

import (
	"github.com/enetx/g"
)

func main() {
	a := g.SliceOf(1, 2, 3)
	println(a.Iter().Count()) // 3

	a = g.SliceOf(1, 2, 3, 4, 5)
	println(a.Iter().Count()) // 5

	a = g.SliceOf(1, 2, 4, 3, 4, 5, 4, 1, 1, 4, 4)
	println(a.Iter().Filter(func(i int) bool { return i == 4 }).Count()) // 5
}
