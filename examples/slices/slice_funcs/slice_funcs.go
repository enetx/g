package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// convert std slice of strings to g.Slice of g.String
	strslice := []string{"aa", "bb", "cc"}
	gs := Slice[string](strslice).Iter().Map(NewString).Collect() // g.Slice[g.String]

	gs.SortBy(func(a, b String) cmp.Ordering { return a.Cmp(b) })
	fmt.Println(gs)

	// convert std slice of ints to g.Slice of g.Int
	intslice := []int{1, 2, 3}
	Slice[int](intslice).Iter().Map(NewInt).Collect() // g.Slice[g.Int]

	// convert std slice of floats to g.Slice of g.Float
	floatslice := []float64{1.1, 2.2, 3.3}
	Slice[float64](floatslice).Iter().Map(NewFloat).Collect() // g.Slice[g.Float]

	////////////////////////////////////////////////////////////////////////////

	s := Slice[Int]{1, 2, 3, 4, 5}

	ss := s.Iter().Map(Int.String).Collect() // g.Slice[g.String]

	ss.Get(0).Some().Format("hello {}").Println()

	is := ss.Iter().Map(String.TryInt).Map(Result[Int].Unwrap).Collect() // g.Slice[g.Int]
	is.Get(0).Some().Add(99).Println()

	////////////////////////////////////////////////////////////////////////////

	ss1 := SliceOf[String]("1", "22", "3a", "44")

	// is1 := ss1.Iter().Map(String.TryInt).Filter(Result[Int].IsOk).Collect()
	// is1.Iter().Map(Result[Int].Ok).Collect().Println() // Slice[1, 22, 44]

	SeqResult[Int](ss1.Iter().Map(String.TryInt)).Ok().Collect() // Slice[1, 22, 44]
}
