package main

import . "github.com/enetx/g"

func main() {
	// Create a new ordered map with integer keys and string values
	m := NewMapOrd[int, string]()
	m.Set(0, "aa")
	m.Set(1, "bb")
	m.Set(2, "cc")
	m.Set(3, "ee")
	m.Set(4, "ff")
	m.Set(5, "gg")
	m.Set(6, "aa")
	m.Iter().
		StepBy(2). // Iterate over the map with a step size of 2
		Collect().
		Println() // MapOrd{0:aa, 2:cc, 4:ff, 6:aa}
}
