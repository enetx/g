package main

import (
	"github.com/enetx/g"
)

func main() {
	g.SliceOf[g.Int](1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		Find(func(v g.Int) bool { return v%2 == 0 }).
		UnwrapOrDefault().
		Print() // 4

	m := g.NewMap[g.Int, g.Int]().Set(1, 11).Set(2, 22).Set(3, 33)
	m.
		Iter().
		Find(func(_, v g.Int) bool { return v == 22 }).
		UnwrapOrDefault().
		Key.
		Print() // 2

	g.MapOrdFromMap(m).
		Iter().
		Find(func(_, v g.Int) bool { return v == 33 }).
		UnwrapOrDefault().
		Key.
		Print() // 3
}
