package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	// Create two slices of integers
	slice1 := g.Slice[int]{1, 2, 3}.Iter()
	slice2 := g.Slice[int]{4, 5, 6}.Iter()

	// Zip the slices in parallel and collect the resulting tuples
	zipped := slice1.Zip(slice2).Collect()
	zipped.Print()

	// Iterate over the zipped tuples and print each one
	for _, v := range zipped {
		fmt.Println(v)
	}
}
