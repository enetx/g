package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	String("foo\r\nbar\n\nbaz\n").
		Lines().
		Exclude(f.IsZero).
		Collect().
		Println() // Slice[foo, bar, baz]

	s := NewString("💛💚💙💜")

	s.LeftJustify(10, "*").Println()  // 💛💚💙💜******
	s.RightJustify(10, "*").Println() // ******💛💚💙💜
	s.Center(10, "*").Println()       // ***💛💚💙💜***

	///////////////////////////////////////////////////////////////////////

	ss := String("Hello, [world]! How [are] you?")
	cuted := NewSlice[String]()
	for ss.ContainsAll("[", "]") {
		var cut String
		ss, cut = ss.Cut("[", "]", true)
		cuted.Push(cut)
	}

	cuted.Println()

	ss.Println()

	println(ss.Contains("Hello"))

	NewString(byte('g')).Println()
	NewString(rune('g')).Println()
	NewString([]rune("hello")).Println()
	NewString([]byte("hello")).Println()

	NewString("").Random(10).Println()
	NewString("").Random(10, ASCII_LETTERS).Println()
	NewString("").Random(10, DIGITS).Println()
	NewString("").Random(10, PUNCTUATION).Println()

	String("https://www.test.com/?query=Hellö Wörld&param=value").
		Encode().
		URL().
		Println() // https://www.test.com/?query=Hell%C3%B6+W%C3%B6rld&param=value

	String("Hellö Wörld@Golang").Encode().URL().Println()   // Hell%C3%B6+W%C3%B6rld@Golang
	String("Hellö Wörld@Golang").Encode().URL("").Println() // Hell%C3%B6+W%C3%B6rld%40Golang

	original := String("Hello, world! This is a test.")
	modified := original.Remove(
		"Hello",
		"test",
	)

	modified.Println()

	num := String("hello title")

	num.Transform(String.Title).Println() // String type

	String("a1b2c3d4e5").
		Chars().
		Filter(String.IsDigit).
		Collect().
		Join().
		Println()

	String("Hello, World!").Truncate(5).Println()
	// result2: "Hello..."

	String("Short").Truncate(10).Println()
	// result2: "Short"

	String("😊😊😊😊😊").Truncate(3).Println()
	// result3: "😊😊😊..."
}
