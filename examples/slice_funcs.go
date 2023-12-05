package main

import "gitlab.com/x0xO/g"

func main() {
	// convert std slice of strings to g.Slice of g.String
	strslice := []string{"aa", "bb", "cc"}
	g.SliceMap(strslice, g.NewString) // g.Slice[g.String]

	// convert std slice of ints to g.Slice of g.Int
	intslice := []int{1, 2, 3}
	g.SliceMap(intslice, g.NewInt) // g.Slice[g.Int]

	// convert std slice of floats to g.Slice of g.Float
	floatslice := []float64{1.1, 2.2, 3.3}
	g.SliceMap(floatslice, g.NewFloat) // g.Slice[g.Float]

	////////////////////////////////////////////////////////////////////////////

	s := g.SliceOf[g.Int](1, 2, 3, 4, 5)

	ss := g.SliceMap(s, g.Int.ToString) // g.Slice[g.String]
	ss.Get(0).Format("hello %s").Print()

	is := g.SliceMap(ss, func(s g.String) g.Int { return s.ToInt().Unwrap() }) // g.Slice[g.Int]
	is.Get(0).Add(99).Print()

	////////////////////////////////////////////////////////////////////////////

	ss1 := g.SliceOf[g.String]("1", "22", "3a", "44")
	is1 := g.SliceMap(ss1, g.String.ToInt).Filter(g.Result[g.Int].IsOk)

	g.SliceMap(is1, g.Result[g.Int].Ok).Print() // Slice[1, 22, 44]
}
