package main

import (
	"github.com/enetx/g/fs"
)

func main() {
	// Copy the contents of the current directory to a new directory named "copy".
	fs.NewDir(".").Copy("copy").Unwrap()

	// Copy the contents of the current directory to a new directory named "copy" while ignoring symbolic links.
	fs.NewDir(".").Copy("copy", false).Unwrap()
}
