package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	sl := g.SliceOf(1, 2, 3, 4, 4, 2, 5, 0)
	s := g.SetOf(sl...) // convert Slice to Set

	s.Iter().ForEach(func(i int) {
		fmt.Println(i)
	})

	even := s.Iter().Filter(func(val int) bool { return val%2 == 0 }).Collect()

	s.Iter().
		Filter(func(val int) bool { return val%2 == 0 }).
		Exclude(f.Zero).
		Inspect(func(i int) {
			fmt.Println(i)
		}).
		Collect().
		Print()

	s2 := g.SetOf(4, 5, 6, 7, 8)
	s.SymmetricDifference(s2).Collect().Print()

	set5 := g.SetOf(1, 2)
	set6 := g.SetOf(2, 3, 4, 9)

	s.Iter().Chain(set5.Iter(), set6.Iter()).Map(func(i int) int { return i + i }).Collect().Print()

	set7 := set5.Difference(set6).Collect()
	set7.Print()

	s = g.SetOf(1, 2, 3, 4, 5)
	even = s.Iter().Filter(f.Even).Collect()
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
