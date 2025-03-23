package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a new ordered map with integer keys and string values
	mo := NewMapOrd[int, string]()
	mo.Set(0, "aa")
	mo.Set(1, "bb")
	mo.Set(2, "cc")
	mo.Set(3, "dd")
	mo.Set(4, "ee")
	mo.Set(5, "ff")
	mo.Set(6, "gg")
	mo.Iter().
		StepBy(2).                                                        // Iterate with a step size of 2
		Exclude(func(k int, _ string) bool { return k == 4 }).            // Exclude entry with key 4
		Inspect(func(k int, v string) { fmt.Println("~inspect", k, v) }). // Inspect each entry
		Map(func(k int, v string) (int, string) { return k, v + v }).
		// Map values to concatenate them with themselves
		Collect(). // Collect the resulting ordered map
		Println()  // Print the ordered map
}
