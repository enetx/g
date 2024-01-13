package main

import "gitlab.com/x0xO/g"

func main() {
	// Create three slices of integers
	slice1 := g.Slice[int]{1, 2, 3}.Iter()
	slice2 := g.Slice[int]{4, 5, 6}.Iter()
	slice3 := g.Slice[int]{7, 8, 9}.Iter()

	// Zip the three slices in parallel and collect the resulting tuples
	zipped := slice1.Zip(slice2, slice3).Collect()

	// Iterate over the zipped tuples and print each one
	for _, v := range zipped {
		v.Print()
	}
}
