package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a new ordered map with integer keys and string values
	g.NewMapOrd[int, string]().
		Set(0, "aa").
		Set(1, "bb").
		Set(2, "cc").
		Set(3, "dd").
		Set(4, "ee").
		Set(5, "ff").
		Set(6, "gg").
		Iter().
		StepBy(2). // Iterate over the map with a step size of 2
		Collect().
		Print()
}
