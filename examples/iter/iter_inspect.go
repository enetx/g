package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	// Create a new ordered map with integer keys and string values
	mo := g.NewMapOrd[int, string]()
	mo.Set(0, "aa").
		Set(1, "bb").
		Set(2, "cc").
		Set(3, "dd").
		Set(4, "ee").
		Set(5, "ff").
		Set(6, "gg").
		Iter().
		StepBy(2).                                                        // Iterate with a step size of 2
		Exclude(func(k int, _ string) bool { return k == 4 }).            // Exclude entry with key 4
		Inspect(func(k int, v string) { fmt.Println("~inspect", k, v) }). // Inspect each entry
		Map(func(k int, v string) (int, string) { return k, v + v }).
		// Map values to concatenate them with themselves
		Collect(). // Collect the resulting ordered map
		Print()    // Print the ordered map
}
