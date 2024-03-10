package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	a := g.SliceOf(1, 2, 3)
	println(a.Iter().Count()) // 3

	a = g.SliceOf(1, 2, 3, 4, 5)
	println(a.Iter().Count()) // 5
}
