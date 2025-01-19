package main

import . "github.com/enetx/g"

func main() {
	// Create a slice of integers
	is := SliceOf[Int](1, 2, 3, 4, 5)

	// Transform the slice of integers into a slice of strings
	itos := TransformSlice(is, Int.String)

	// Iterate over the transformed slice, perform folding, and print the result
	itos.Iter().
		Fold("0", // Initial accumulator value
			func(acc, val String) String {
				// Folding function: concatenate each element in the iterator with the accumulator
				return Sprintf("({} + {})", acc, val)
				}).
		Println() // (((((0 + 1) + 2) + 3) + 4) + 5)
}
