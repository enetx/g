package main

import (
	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	// Open a new file with the specified name "text.txt"
	g.NewFile("text.txt").
		Lines().                 // Read the file line by line
		Unwrap().                // Unwrap the Result type to get the underlying iterator
		Skip(3).                 // Skip the first 3 lines
		Exclude(filters.IsZero). // Exclude lines that are empty or contain only whitespaces
		Dedup().                 // Remove consecutive duplicate lines
		Map(g.String.Upper).     // Convert each line to uppercase
		ForEach(                 // For each line, print it
			func(s g.String) {
				s.Print()
			})
}
