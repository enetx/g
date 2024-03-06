package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a new ordered map with integer keys and string values
	m := g.NewMapOrd[int, string]()
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
		Print()

	g.SliceOf(0, 1, 2, 3, 4, 5, 6).
		Iter().
		StepBy(2).
		Collect().
		Print()
}
