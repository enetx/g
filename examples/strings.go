package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	String("foo\r\nbar\n\nbaz\n").
		Lines().
		Exclude(f.Zero).
		Collect().
		Print() // Slice[foo, bar, baz]

	s := NewString("💛💚💙💜")

	s.LeftJustify(10, "*").Print()  // 💛💚💙💜******
	s.RightJustify(10, "*").Print() // ******💛💚💙💜
	s.Center(10, "*").Print()       // ***💛💚💙💜***

	///////////////////////////////////////////////////////////////////////

	ss := String("Hello, [world]! How [are] you?")

	cuted := NewSlice[String]()

	for ss.ContainsAll("[", "]") {
		var cut String
		ss, cut = ss.Cut("[", "]")
		cuted.AppendInPlace(cut)
	}

	cuted.Print()

	NewString(byte('g')).Print()
	NewString(rune('g')).Print()
	NewString([]rune("hello")).Print()
	NewString([]byte("hello")).Print()

	NewString("").Random(10).Print()
	NewString("").Random(10, ASCII_LETTERS).Print()
	NewString("").Random(10, DIGITS).Print()
	NewString("").Random(10, PUNCTUATION).Print()

	String("https://www.test.com/?query=Hellö Wörld&param=value").
		Encode().
		URL().
		Print() // https://www.test.com/?query=Hell%C3%B6+W%C3%B6rld&param=value

	String("Hellö Wörld@Golang").Encode().URL().Print()   // Hell%C3%B6+W%C3%B6rld@Golang
	String("Hellö Wörld@Golang").Encode().URL("").Print() // Hell%C3%B6+W%C3%B6rld%40Golang

	original := String("Hello, world! This is a test.")
	modified := original.Remove(
		"Hello",
		"test",
	)

	modified.Print()
}
