package main

import "gitlab.com/x0xO/g"

func main() {
	s := g.NewString("💛💚💙💜")

	s.LeftJustify(10, "*").Print()
	s.RightJustify(10, "*").Print()
	s.Center(10, "*").Print()

	// 💛💚💙💜******
	// ******💛💚💙💜
	// ***💛💚💙💜***

	///////////////////////////////////////////////////////////////////////

	ss := g.String("Hello, [world]! How [are] you?")

	cuted := g.NewSlice[g.String]()

	for ss.ContainsAll("[", "]") {
		var cut g.String
		ss, cut = ss.Cut("[", "]")
		cuted.AppendInPlace(cut)
	}

	cuted.Print()

	g.NewString(byte('g')).Print()
	g.NewString(rune('g')).Print()
	g.NewString([]rune("hello")).Print()
	g.NewString([]byte("hello")).Print()

	g.NewString("").Random(10).Print()
	g.NewString("").Random(10, g.ASCII_LETTERS).Print()
	g.NewString("").Random(10, g.DIGITS).Print()
	g.NewString("").Random(10, g.PUNCTUATION).Print()

	g.NewString("some text").Comp().Brotli().Decomp().Brotli().Unwrap().Print()
}
