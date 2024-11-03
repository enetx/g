package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	evens, odds := SliceOf(1, 2, 3, 4, 5).
		Iter().
		Partition(f.Even)

	fmt.Println("Even numbers:", evens) // Output: Even numbers: Slice[2, 4]
	fmt.Println("Odd numbers:", odds)   // Output: Odd numbers: Slice[1, 3, 5]
}
