package main

import (
	"fmt"
	"io"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	// Example 1: Reading and processing lines

	// Open a new file with the specified name "text.txt" and process it line by line.
	NewFile("text.txt").
		Lines().                            // Reads the file line by line.
		Skip(2).                            // Skips the first 2 lines in the iterator.
		Exclude(f.IsZero).                  // Excludes lines that are empty or contain only whitespaces.
		Dedup().                            // Removes consecutive duplicate lines.
		Map(String.Upper).                  // Converts each line to uppercase.
		Range(func(s Result[String]) bool { // Iterates over the lines while a condition is true.
			if s.IsErr() { // Handles any errors encountered while reading lines.
				fmt.Println("Error:", s.Err())
				return false // Stops the iteration if an error occurs.
			}

			if s.V().Contains("COULD") { // Checks if the line contains the substring "COULD".
				return false // Stops the iteration if the condition is met.
			}

			fmt.Println(s.V()) // Prints the line.
			return true        // Continues the iteration.
		})

	// Example 2: Reading and processing chunks

	offset := int64(10) // Initialize offset.

	// Open the file and read it in chunks of 3 bytes.
	NewFile("text.txt").
		Seek(offset, io.SeekStart).Unwrap(). // Moves the read pointer to 10 bytes from the beginning of the file.
		Chunks(3).                           // Reads the file in chunks of 3 bytes.
		Inspect(func(s String) {             // Inspects each chunk during iteration.
			// Updates the offset based on the length of each chunk.
			// This keeps track of the current position in the file.
			offset += int64(s.Bytes().Len())
		}).
		Collect(). // Collects all chunks into a single collection (Slice).
		Unwrap().  // Unwraps the result, extracting the successful value or panicking on error.
		Join().    // Joins all collected chunks into a single string.
		Println()  // Prints the joined string.

	fmt.Println(offset) // Prints the final offset after processing the chunks.
}
