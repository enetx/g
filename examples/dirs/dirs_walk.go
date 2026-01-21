package main

import . "github.com/enetx/g"

func main() {
	// Recursively walk through the directory tree starting from the current directory
	NewDir(".").Walk().
		// Exclude directories and symlinked directories
		Exclude(func(f *File) bool { return f.IsDir() && f.Dir().Ok().IsLink() }).
		// Exclude all symbolic links (files or directories)
		Exclude((*File).IsLink).
		// Process each walk result
		ForEach(func(v Result[*File]) {
			if v.IsOk() {
				// Print the path of the file if no error occurred
				v.Ok().Path().Ok().Println()
			}
		})
}
