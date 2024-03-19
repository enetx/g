package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	evens, odds := g.SliceOf(1, 2, 3, 4, 5).
		Iter().
		Partition(
			func(v int) bool {
				return v%2 == 0
			})

	fmt.Println("Even numbers:", evens) // Output: Even numbers: Slice[2, 4]
	fmt.Println("Odd numbers:", odds)   // Output: Odd numbers: Slice[1, 3, 5]
}
