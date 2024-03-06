package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	iter := g.Slice[int]{1, 2, 3, 4, 5, 6}.
		Iter().
		Skip(3).
		Collect()

	fmt.Println(iter)
}
