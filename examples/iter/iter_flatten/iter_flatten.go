package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Flatten a slice containing various types of elements
	Slice[any]{
		1,                            // integer
		SliceOf(2, 3),                // slice of integers
		"abc",                        // string
		SliceOf("awe", "som", "e"),   // slice of strings
		SliceOf("co", "ol"),          // another slice of strings
		SliceOf(4.5, 6.7),            // slice of floats
		map[string]string{"a": "ss"}, // map with string keys and values
		SliceOf(
			Map[int, int]{1: 1}, // slice of maps
			Map[int, int]{2: 2},
		),
	}.
		Iter().    // creates an iterator for subsequent operations
		Flatten(). // flattens nested slices, transforming them into a flat slice
		Collect(). // gathers the elements of the iterator into a new slice.
		Println()
		// outputs the elements of the slice to the console. Slice[1, 2, 3, abc, awe, som, e, co, ol, 4.5, 6.7, map[a:ss], Map{1:1}, Map{2:2}]

	// Example 2: Flatten a slice of strings by individual characters
	words := SliceOf[String]("alpha", "beta", "gamma", "💛💚💙💜", "世界")

	// FlatMap maps each string to its characters and flattens the result —
	// fully typed, no any round-trips (generic methods, Go 1.27).
	words.Iter().
		FlatMap(String.Chars).
		Map(String.Upper).
		Collect().
		Join().
		Println() // ALPHABETAGAMMA💛💚💙💜世界

	// Example 3: Check if the flattened slice contains a specific element.
	// FlatMap with a method expression keeps the element type: Slice[string], not Slice[any].
	ch := Slice[Slice[string]]{{"a", "b", "c"}, {"d", "f", "g"}}.
		Iter().
		FlatMap(Slice[string].Iter).
		Collect() // Slice[string]{"a", "b", "c", "d", "f", "g"}

	fmt.Println(ch.Contains("x"))         // false
	fmt.Println(ch.Contains("a"))         // true
	fmt.Println(ch.ContainsBy(f.Eq("c"))) // true
}
