package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Flatten a slice containing various types of elements
	g.Slice[any]{
		1,                            // integer
		g.SliceOf(2, 3),              // slice of integers
		"abc",                        // string
		g.SliceOf("awe", "som", "e"), // slice of strings
		g.SliceOf("co", "ol"),        // another slice of strings
		g.SliceOf(4.5, 6.7),          // slice of floats
		map[string]string{"a": "ss"}, // map with string keys and values
		g.SliceOf(
			g.MapOrd[int, int]{{1, 1}}, // slice of ordered maps
			g.MapOrd[int, int]{{2, 2}}),
	}.
		Iter().    // creates an iterator for subsequent operations
		Flatten(). // flattens nested slices, transforming them into a flat slice
		Collect(). // gathers the elements of the iterator into a new slice.
		Print()
		// outputs the elements of the slice to the console. Slice[1, 2, 3, abc, awe, som, e, co, ol, 4.5, 6.7, map[a:ss], {1 1}, {2 2}]

	// Example 2: Flatten a slice of strings by individual characters
	words := g.SliceOf[g.String]("alpha", "beta", "gamma", "💛💚💙💜", "世界")

	// MapSlice applies a mapping function to each element of the source slice and returns a new slice.
	// In this example, it maps each string in 'words' to its individual characters.
	g.SliceMap(words, g.String.Chars).
		AsAny(). // Required if the source slice is not of type g.Slice[any]
		Iter().
		Flatten().
		Collect().
		Join().
		Print() // alphabetagamma💛💚💙💜世界

	// Example 3: Check if the flattened slice contains a specific element
	ch := g.Slice[g.Slice[string]]{{"a", "b", "c"}, {"d", "f", "g"}}.
		AsAny(). // g.Slice[any]{g.Slice[string]{"a", "b", "c"}, g.Slice[string]{"d", "f", "g"}}
		Iter().
		Flatten().
		Collect() // g.Slice[any]{"a", "b", "c", "d", "f", "g"}

	fmt.Println(ch.Contains("x"))              // false
	fmt.Println(ch.Contains("a"))              // true
	fmt.Println(ch.Contains(4444))             // false
	fmt.Println(ch.ContainsBy(f.Eq[any]("c"))) // true
}
