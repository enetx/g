package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a slice of integers
	is := g.SliceOf[g.Int](1, 2, 3, 4, 5)

	// Transform the slice of integers into a slice of strings
	itos := g.TransformSlice(is, g.Int.ToString)

	// Iterate over the transformed slice, perform folding, and print the result
	itos.Iter().
		Fold("0", // Initial accumulator value
			func(acc, val g.String) g.String {
				// Folding function: concatenate each element in the iterator with the accumulator
				return g.Sprintf("(%s + %s)", acc, val)
			}).
		Print()
}