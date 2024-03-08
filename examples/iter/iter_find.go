package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	r := g.SliceOf[g.Int](1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		Find(func(v g.Int) bool { return v%2 == 0 })

	r.UnwrapOr(10).Print()
}
