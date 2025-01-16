package main

import . "github.com/enetx/g"

func main() {
	NewDir(".").Walk(walker)
}

func walker(f *File) error {
	// Skip symbolic link directories
	if f.IsDir() && f.Dir().Ok().IsLink() {
		return SkipWalk
	}

	// Skip symbolic link files
	if f.IsLink() {
		return nil
	}

	// Print the path
	f.Path().Ok().Println()

	return nil
}
