package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	sl := SliceOf(1, 2, 3, 4, 4, 2, 5, 0)
	s := SetOf(sl...) // convert Slice to Set

	s.Iter().ForEach(func(i int) {
		fmt.Println(i)
	})

	even := s.Iter().Filter(func(val int) bool { return val%2 == 0 }).Collect()

	s.Iter().
		Filter(func(val int) bool { return val%2 == 0 }).
		Exclude(f.IsZero).
		Inspect(func(i int) {
			fmt.Println(i)
		}).
		Collect().
		Println()

	s2 := SetOf(4, 5, 6, 7, 8)
	s.SymmetricDifference(s2).Collect().Println()

	set5 := SetOf(1, 2)
	set6 := SetOf(2, 3, 4, 9)

	s.Iter().Chain(set5.Iter(), set6.Iter()).Map(func(i int) int { return i + i }).Collect().Println()

	set7 := set5.Difference(set6).Collect()
	set7.Println()

	s = SetOf(1, 2, 3, 4, 5)
	even = s.Iter().Filter(f.IsEven).Collect()
	even.Println()

	s = s.Remove(1)
	s.Println()

	// iterate over set
	for value := range s {
		fmt.Println(value)
	}

	s = SetOf(1, 2, 3)
	TransformSet(s, NewInt) // Set[g.Int]
}
