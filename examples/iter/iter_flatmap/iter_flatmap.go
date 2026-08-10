package main

import . "github.com/enetx/g"

func main() {
	// Example 1: FlatMap to expand numbers into ranges
	SliceOf(1, 2, 3).
		Iter().
		FlatMap(func(n int) Seq[int] {
			// Expand each number to [n, n*10, n*100]
			return SliceOf(n, n*10, n*100).Iter()
		}).
		Collect().
		Slice().
		Println() // Slice[1, 10, 100, 2, 20, 200, 3, 30, 300]

	// Example 2: FlatMap to split sentences into words
	SliceOf[String]("hello world", "go programming").
		Iter().
		// Split each sentence into words
		FlatMap(func(s String) Seq[String] { return SliceOf(s.Fields()...).Iter() }).
		Collect().
		Slice().
		Println() // Slice[hello, world, go, programming]

	// Example 3: FlatMap with conditional expansion
	SliceOf(1, 2, 3, 4, 5).
		Iter().
		FlatMap(func(n int) Seq[int] {
			// Even numbers expand to [n, n*2], odd numbers to just [n]
			if n%2 == 0 {
				return SliceOf(n, n*2).Iter()
			}
			return SliceOf(n).Iter()
		}).
		Collect().
		Slice().
		Println() // Slice[1, 2, 4, 3, 4, 8, 5]

	// Example 4: FlatMap to generate coordinates
	SliceOf[String]("A", "B").
		Iter().
		FlatMap(func(letter String) Seq[String] {
			// For each letter, generate coordinates with numbers
			coords := NewSlice[String]()
			for i := range Int(3) {
				coords.Push(letter.Append(i.String()))
			}
			return coords.Iter()
		}).
		Collect().
		Slice().
		Println() // Slice[A0, A1, A2, B0, B1, B2]

	// FlatMap is type-changing too: each String expands into its characters
	SliceOf[String]("go", "1.27").
		Iter().
		FlatMap(func(s String) Seq[String] { return SliceOf(s.Chars()...).Iter() }).
		Map(String.Upper).
		Collect().
		Slice().
		Println() // Slice[G, O, 1, ., 2, 7]

}
