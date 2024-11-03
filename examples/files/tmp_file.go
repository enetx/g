package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a temporary file in the current directory with a ".txt" extension
	f := NewFile("").CreateTemp("./", "*.txt").Unwrap().Write("some text").Unwrap()

	// Alternatively, create a temporary file without specifying the extension
	// f := g.NewFile("").CreateTemp().Unwrap().Write("some text").Unwrap()

	// Print the path and content of the temporary file
	fmt.Println(f.Path().Unwrap(), f.Read().Unwrap())

	// Calculate the MD5 hash of the file's content and print it
	f.Read().Unwrap().Hash().MD5().Print()

	// Remove the temporary file and print whether it still exists
	fmt.Println(f.Remove().Unwrap().Exist())
}
