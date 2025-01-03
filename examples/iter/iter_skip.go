package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	iter := Slice[int]{1, 2, 3, 4, 5, 6}.
		Iter().
		Skip(3).
		Collect()

	fmt.Println(iter) // Slice[4, 5, 6]
}
