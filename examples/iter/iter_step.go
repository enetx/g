package main

import . "github.com/enetx/g"

func main() {
	// Create a new ordered map with integer keys and string values
	m := NewMapOrd[int, string]()
	m.
		Set(0, "aa").
		Set(1, "bb").
		Set(2, "cc").
		Set(3, "ee").
		Set(4, "ff").
		Set(5, "gg").
		Set(6, "aa").
		Iter().
		StepBy(2). // Iterate over the map with a step size of 2
		Collect().
		Println() // MapOrd{0:aa, 2:cc, 4:ff, 6:aa}
}
