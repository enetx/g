package main

import (
	. "github.com/enetx/g"

	"github.com/enetx/g/fs"
)

func main() {
	// Recursively walk through the directory tree starting from the current directory
	fs.NewDir(".").Walk().
		// Exclude directories and symlinked directories
		Exclude(func(f *fs.File) bool { return f.IsDir() && f.Dir().Ok().IsLink() }).
		// Exclude file symlinks
		Exclude((*fs.File).IsLink).
		// Process each walk result
		ForEach(func(v Result[*fs.File]) {
			if v.IsOk() {
				// Print the path of the file if no error occurred
				v.Ok().Path().Ok().Println()
			}
		})
}
