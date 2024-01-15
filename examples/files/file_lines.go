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
		Seek(10).Ok().           // Seek to 10 bytes from the beginning of the file
		Skip(3).                 // Skip the first 3 lines
		Exclude(filters.IsZero). // Exclude lines that are empty or contain only whitespaces
		Dedup().                 // Remove consecutive duplicate lines
		Map(g.String.Upper).     // Convert each line to uppercase
		Range(func(s g.String) bool {
			// Check if the line contains "COULD"
			if s.Contains("COULD") {
				f.Close() // Close the file if "COULD" is found
				return false
			}

			fmt.Println(s)
			return true
		})

	// Example 2: Reading and processing chunks

	offset := int64(10) // Initialize offset

	g.NewFile("text.txt").
		Chunks(20).        // Read the file in chunks of 3 bytes
		Unwrap().          // Unwrap the Result type to get the underlying iterator
		Seek(offset).Ok(). // Seek to 10 bytes from the beginning of the file
		Inspect(func(s g.String) {
			offset += int64(s.ToBytes().Len()) // Update the offset based on the length of each chunk
		}).
		ForEach(func(s g.String) {
			fmt.Print(s) // Print each chunk
		})

	// Print the final offset
	fmt.Println(offset)
}
