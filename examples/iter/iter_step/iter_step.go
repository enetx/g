package main

import . "github.com/enetx/g"

func main() {
	// Create a new ordered map with integer keys and string values
	m := NewMapOrd[int, string]()
	m.Insert(0, "aa")
	m.Insert(1, "bb")
	m.Insert(2, "cc")
	m.Insert(3, "ee")
	m.Insert(4, "ff")
	m.Insert(5, "gg")
	m.Insert(6, "aa")
	m.Iter().
		StepBy(2). // Iterate over the map with a step size of 2
		Collect().
		Println() // MapOrd{0:aa, 2:cc, 4:ff, 6:aa}
}
