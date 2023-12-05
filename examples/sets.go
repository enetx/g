package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	sl := g.SliceOf(1, 2, 3, 4, 4, 2, 5)
	s := g.SetOf(sl...) // convert Slice to Set
	s.Print()

	s2 := g.SetOf(4, 5, 6, 7, 8)
	s.SymmetricDifference(s2).Print()

	set5 := g.SetOf(1, 2)
	set6 := g.SetOf(2, 3, 4)

	set7 := set5.Difference(set6)
	set7.Print()

	s = g.SetOf(1, 2, 3, 4, 5)
	even := s.Filter(func(val int) bool { return val%2 == 0 })
	even.Print()

	s = s.Remove(1)
	s.Print()

	// iterate over set
	for value := range s {
		fmt.Println(value)
	}

	s = g.SetOf(1, 2, 3)
	g.SetMap(s, g.NewInt) // g.Set[g.Int]
}
