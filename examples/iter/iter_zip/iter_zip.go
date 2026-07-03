package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create two slices of integers
	slice1 := Slice[int]{1, 2, 3}.Iter()
	slice2 := Slice[int]{4, 5, 6}.Iter()

	// Zip the slices in parallel and collect the resulting pairs
	zipped := slice1.Zip(slice2).Collect()
	fmt.Println(zipped) // [{1 4} {2 5} {3 6}]

	// Iterate over the zipped pairs and print each one
	for _, p := range zipped {
		fmt.Println(p.Key, p.Value)
	}

	// Zip pairs sequences of DIFFERENT types; SeqPairs re-enters the main
	// chain via Map, or splits back into two typed slices via Unzip
	ids := SliceOf(1, 2, 3)
	names := SliceOf[String]("alice", "bob", "carol")

	ids.Iter().
		Zip(names.Iter()).
		Map(func(id int, name String) String { return Format("{}:{}", id, name) }).
		Collect().
		Join(" ").
		Println() // 1:alice 2:bob 3:carol

	keys, values := ids.Iter().
		Zip(names.Iter()).
		Filter(func(id int, _ String) bool { return id%2 == 1 }).
		Unzip()

	keys.Println()   // Slice[1, 3]
	values.Println() // Slice[alice, carol]

	if p := ids.Iter().Zip(names.Iter()).Find(func(_ int, n String) bool { return n == "bob" }); p.IsSome() {
		fmt.Println(p.Some().Key, p.Some().Value) // 2 bob
	}
}
