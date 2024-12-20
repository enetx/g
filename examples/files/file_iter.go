package main

import (
	"fmt"
	"io"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Reading and processing lines

	// Open a new file with the specified name "text.txt" and read the file line by line
	r := NewFile("text.txt").Lines()

	// switch pattern
	switch {
	case r.IsOk():
		r.Ok().
			Skip(3).           // Skip the first 3 lines
			Exclude(f.IsZero). // Exclude lines that are empty or contain only whitespaces
			Dedup().           // Remove consecutive duplicate lines
			Map(String.Upper). // Convert each line to uppercase
			Range(func(s String) bool {
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

	NewFile("text.txt").
		Seek(offset, io.SeekStart).Ok(). // Seek to 10 bytes from the beginning of the file
		Chunks(3).                       // Read the file in chunks of 3 bytes
		Unwrap().                        // Unwrap the Result type to get the underlying iterator
		Inspect(func(s String) {         // Update the offset based on the length of each chunk
			offset += int64(s.Bytes().Len())
		}).
		ForEach(func(s String) { // For each chunk, print it
			fmt.Print(s)
		})

	// Print the final offset
	fmt.Println(offset)
}
