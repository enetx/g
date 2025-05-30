package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// convert std slice of strings to g.Slice of g.String
	strslice := []string{"aa", "bb", "cc"}
	gs := TransformSlice(strslice, NewString) // g.Slice[g.String]

	gs.SortBy(func(a, b String) cmp.Ordering { return a.Cmp(b) })
	fmt.Println(gs)

	// convert std slice of ints to g.Slice of g.Int
	intslice := []int{1, 2, 3}
	TransformSlice(intslice, NewInt) // g.Slice[g.Int]

	// convert std slice of floats to g.Slice of g.Float
	floatslice := []float64{1.1, 2.2, 3.3}
	TransformSlice(floatslice, NewFloat) // g.Slice[g.Float]

	////////////////////////////////////////////////////////////////////////////

	s := SliceOf[Int](1, 2, 3, 4, 5)

	ss := TransformSlice(s, Int.String) // g.Slice[g.String]
	ss.Get(0).Some().Format("hello %s").Println()

	is := TransformSlice(ss, func(s String) Int { return s.ToInt().Unwrap() }) // g.Slice[g.Int]
	is.Get(0).Some().Add(99).Println()

	////////////////////////////////////////////////////////////////////////////

	ss1 := SliceOf[String]("1", "22", "3a", "44")
	is1 := TransformSlice(ss1, String.ToInt).Iter().Filter(Result[Int].IsOk).Collect()

	TransformSlice(is1, Result[Int].Ok).Println() // Slice[1, 22, 44]
}
