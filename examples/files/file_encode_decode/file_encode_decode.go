package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Encoding and decoding using Gob

	// Encode a slice of integers to a Gob file
	NewFile("test.gob").Encode().Gob(SliceOf(1, 2, 3)).Unwrap()

	// Decode the Gob file into a new slice of integers
	var gobdata Slice[int]
	NewFile("test.gob").Decode().Gob(&gobdata)

	// Print the decoded Gob data
	fmt.Println(gobdata)

	// Encoding and decoding using JSON

	// Encode a slice of integers to a JSON file
	NewFile("test.json").Encode().JSON(SliceOf(1, 2, 3)).Unwrap()

	// Decode the JSON file into a new slice of integers
	var jsondata Slice[int]
	NewFile("test.json").Decode().JSON(&jsondata)

	// Print the decoded JSON data
	fmt.Println(jsondata)
}
