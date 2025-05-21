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
			MapOrd[int, int]{{1, 1}}, // slice of ordered maps
			MapOrd[int, int]{{2, 2}}),
	}.
		Iter().    // creates an iterator for subsequent operations
		Flatten(). // flattens nested slices, transforming them into a flat slice
		Collect(). // gathers the elements of the iterator into a new slice.
		Println()
		// outputs the elements of the slice to the console. Slice[1, 2, 3, abc, awe, som, e, co, ol, 4.5, 6.7, map[a:ss], {1 1}, {2 2}]

	// Example 2: Flatten a slice of strings by individual characters
	words := SliceOf[String]("alpha", "beta", "gamma", "💛💚💙💜", "世界")

	// MapSlice applies a mapping function to each element of the source slice and returns a new slice.
	// In this example, it maps each string in 'words' to its individual characters.
	TransformSlice(words, func(w String) Slice[String] { return w.Chars().Collect() }).
		// SliceMap(words, String.Chars).
		AsAny(). // Required if the source slice is not of type g.Slice[any]
		Iter().
		Map(func(a any) any { return a.(Slice[String]).Iter().Map(String.Upper).Collect().AsAny() }).
		Flatten().
		Collect().
		Join().
		Println() // ALPHABETAGAMMA💛💚💙💜世界

	// Example 3: Check if the flattened slice contains a specific element
	ch := Slice[Slice[string]]{{"a", "b", "c"}, {"d", "f", "g"}}.
		AsAny(). // Slice[any]{Slice[string]{"a", "b", "c"}, Slice[string]{"d", "f", "g"}}
		Iter().
		Flatten().
		Collect() // Slice[any]{"a", "b", "c", "d", "f", "g"}

	fmt.Println(ch.Contains("x"))              // false
	fmt.Println(ch.Contains("a"))              // true
	fmt.Println(ch.Contains(4444))             // false
	fmt.Println(ch.ContainsBy(f.Eq[any]("c"))) // true
}
