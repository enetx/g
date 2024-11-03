package main

import . "github.com/enetx/g"

func main() {
	// Copy the contents of the current directory to a new directory named "copy".
	NewDir(".").Copy("copy").Unwrap()

	// Copy the contents of the current directory to a new directory named "copy" while ignoring symbolic links.
	NewDir(".").Copy("copy", false).Unwrap()
}
