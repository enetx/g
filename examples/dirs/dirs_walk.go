package main

import "gitlab.com/x0xO/g"

func main() {
	// Create a directory instance for the current directory, read its contents recursively, and iterate over files
	g.NewDir("./").Read(true).Unwrap().Iter().ForEach(walker)
}

// Walker function to process each file in the directory
func walker(f *g.File) {
	// Check if the current file is a directory
	if f.Stat().Unwrap().IsDir() {
		// If it's a directory, recursively walk through its contents and apply the walker function
		f.Dir().Unwrap().Read(true).Unwrap().Iter().ForEach(walker)
	}

	// Print the path of the current file (including both files and directories)
	f.Path().Unwrap().Print()
}
