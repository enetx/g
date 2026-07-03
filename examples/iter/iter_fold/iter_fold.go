package main

import . "github.com/enetx/g"

func main() {
	// Create a slice of integers
	is := SliceOf[Int](1, 2, 3, 4, 5)

	// Transform the slice of integers into a slice of strings
	itos := is.Iter().Map(Int.String).Collect()

	// Iterate over the transformed slice, perform folding, and print the result
	itos.Iter().
		Fold("0", // Initial accumulator value
			func(acc, val String) String {
				// Folding function: concatenate each element in the iterator with the accumulator
				return Format("({} + {})", acc, val)
				}).
		Println() // (((((0 + 1) + 2) + 3) + 4) + 5)

	// Fold accepts an accumulator of any type, not just the element type:
	// fold Slice[String] into an Int
	words := SliceOf[String]("alpha", "beta", "gamma")
	words.Iter().
		Fold(Int(0), func(acc Int, w String) Int { return acc + w.Len() }).
		Println() // 14

	// ...or into a typed counter Map — keys stay String, no any casts
	counts := SliceOf[String]("a", "b", "a").
		Iter().
		Fold(NewMap[String, Int](), func(acc Map[String, Int], s String) Map[String, Int] {
			acc[s]++
			return acc
		})
	counts.Get("a").UnwrapOr(0).Println() // 2

}
