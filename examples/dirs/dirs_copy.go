package main

import "github.com/enetx/g"

func main() {
	// Copy the contents of the current directory to a new directory named "copy".
	g.NewDir(".").Copy("copy").Unwrap()

	// Copy the contents of the current directory to a new directory named "copy" while ignoring symbolic links.
	g.NewDir(".").Copy("copy", false).Unwrap()
}
