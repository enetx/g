package main

import (
	"fmt"
	"io"

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

	// Example 2: Reading and processing chunks

	offset := int64(10) // Initialize offset

	g.NewFile("text.txt").
		Seek(offset, io.SeekStart).Ok(). // Seek to 10 bytes from the beginning of the file
		Chunks(3).                       // Read the file in chunks of 3 bytes
		Unwrap().                        // Unwrap the Result type to get the underlying iterator
		Inspect(func(s g.String) {       // Update the offset based on the length of each chunk
			offset += int64(s.ToBytes().Len())
		}).
		ForEach(func(s g.String) { // For each chunk, print it
			fmt.Print(s)
		})

	// Print the final offset
	fmt.Println(offset)
}
