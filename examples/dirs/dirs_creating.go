package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	// Create a directory instance for the current directory and print its path
	d := g.NewDir("")
	d.Path().Unwrap().Print()

	// Check if the directory exists
	fmt.Println(d.Exist())

	// Create all directories in the specified path, even if some of the intermediate directories don't exist
	g.NewDir("./some/dir/that/dont/exist/").CreateAll().Unwrap()

	// Rename the directory "./some/dir/that/" to "./some/dir/aaa/ccc/" and print the new path
	g.NewDir("./some/dir/that/").Rename("./some/dir/aaa/ccc/").Unwrap().Print()

	// Create all directories in the path "aaa", then rename it to "bbb" and print the new path
	g.NewDir("aaa").CreateAll().Unwrap().Rename("bbb").Unwrap()

	// Create a temporary directory, print its path, and assign it to variable 'd'
	d = g.NewDir("").CreateTemp().Unwrap().Print()

	// Remove the directory and print whether it exists after removal
	d.Remove().Unwrap()
	fmt.Println(d.Exist())
}
