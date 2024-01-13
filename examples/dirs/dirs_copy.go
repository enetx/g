package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a directory instance for the current directory and copy its contents to a new directory named "copy"
	d := g.NewDir(".").Copy("copy").Unwrap()

	// Print the path of the copied directory
	d.Path().Unwrap().Print()
}
