package main

import (
	"github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	g.String("foo\r\nbar\n\nbaz\n").
		Lines().
		Exclude(f.Zero).
		Collect().
		Print() // Slice[foo, bar, baz]

	s := g.NewString("💛💚💙💜")

	s.LeftJustify(10, "*").Print()  // 💛💚💙💜******
	s.RightJustify(10, "*").Print() // ******💛💚💙💜
	s.Center(10, "*").Print()       // ***💛💚💙💜***

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

	g.String("https://www.test.com/?query=Hellö Wörld&param=value").
		Enc().
		URL().
		Print() // https://www.test.com/?query=Hell%C3%B6+W%C3%B6rld&param=value

	g.String("Hellö Wörld@Golang").Enc().URL().Print()   // Hell%C3%B6+W%C3%B6rld@Golang
	g.String("Hellö Wörld@Golang").Enc().URL("").Print() // Hell%C3%B6+W%C3%B6rld%40Golang

	original := g.String("Hello, world! This is a test.")
	modified := original.Remove(
		"Hello",
		"test",
	)

	modified.Print()
}
