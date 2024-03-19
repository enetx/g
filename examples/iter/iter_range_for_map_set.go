package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	// Example 1: Map key-value pairs in a map, print the result, and close the iterator at a specific condition
	m := g.NewMap[int, string]().Set(88, "aa").Set(99, "bb").Set(199, "ii").Iter()

	m.
		Map(func(k int, v string) (int, string) { return k + k, v }).
		Range(func(k int, v string) bool {
			// Close the iterator if the key is 198
			if k == 198 {
				return false
			}

			fmt.Println(k, v)
			return true
		})

	// Example 2: Iterate over a set of integers, print each value, and stop the iteration at a specific condition
	set := g.NewSet[int]().Add(1, 2, 3, 4, 5).Iter()

	set.
		Map(func(v int) int { return v + v }).
		Range(func(v int) bool {
			// Close the iterator if the value is 10
			if v == 10 {
				return false
			}

			fmt.Println(v)
			return true
		})
}
