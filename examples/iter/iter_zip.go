package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create two slices of integers
	slice1 := Slice[int]{1, 2, 3}.Iter()
	slice2 := Slice[int]{4, 5, 6}.Iter()

	// Zip the slices in parallel and collect the resulting tuples
	zipped := slice1.Zip(slice2).Collect()
	zipped.Println() // MapOrd{1:4, 2:5, 3:6}

	// Iterate over the zipped tuples and print each one
	for k, v := range zipped.Iter() {
		fmt.Println(k, v)
	}
}
