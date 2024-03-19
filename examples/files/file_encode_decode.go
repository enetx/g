package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	// Encoding and decoding using Gob

	// Encode a slice of integers to a Gob file
	g.NewFile("test.gob").Enc().Gob(g.SliceOf(1, 2, 3)).Unwrap()

	// Decode the Gob file into a new slice of integers
	var gobdata g.Slice[int]
	g.NewFile("test.gob").Dec().Gob(&gobdata)

	// Print the decoded Gob data
	fmt.Println(gobdata)

	// Encoding and decoding using JSON

	// Encode a slice of integers to a JSON file
	g.NewFile("test.json").Enc().JSON(g.SliceOf(1, 2, 3)).Unwrap()

	// Decode the JSON file into a new slice of integers
	var jsondata g.Slice[int]
	g.NewFile("test.json").Dec().JSON(&jsondata)

	// Print the decoded JSON data
	fmt.Println(jsondata)
}
