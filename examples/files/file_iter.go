package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/filters"
)

func main() {
	// Example 1: Reading and processing lines
	r := g.
		NewFile("text.txt"). // Open a new file with the specified name "text.txt"
		Lines()              // Read the file line by line

	// switch pattern
	switch {
	case r.IsOk():
		r.Ok().
			Skip(3).                 // Skip the first 3 lines
			Exclude(filters.IsZero). // Exclude lines that are empty or contain only whitespaces
			Dedup().                 // Remove consecutive duplicate lines
			Map(g.String.Upper).     // Convert each line to uppercase
			Range(func(s g.String) bool {
				if s.Contains("COULD") { // Check if the line contains "COULD"
					return false
				}

				fmt.Println(s)
				return true
			})
	case r.IsErr():
		fmt.Println("err:", r.Err())
	}
}
