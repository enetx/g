package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	// Example 1: Reading and processing lines
	f := g.NewFile("text.txt") // Open a new file with the specified name "text.txt"
	f.
		Lines().                 // Read the file line by line
		Unwrap().                // Unwrap the Result type to get the underlying iterator
		Skip(3).                 // Skip the first 3 lines
		Exclude(filters.IsZero). // Exclude lines that are empty or contain only whitespaces
		Dedup().                 // Remove consecutive duplicate lines
		Map(g.String.Upper).     // Convert each line to uppercase
		Range(func(s g.String) bool {
			if s.Contains("COULD") { // Check if the line contains "COULD"
				f.Close() // Close the file if "COULD" is found
				return false
			}

			fmt.Println(s)
			return true
		})
}
