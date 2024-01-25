package main

import "gitlab.com/x0xO/g"

func main() {
	g.NewDir(".").Walk(walker)
}

func walker(f *g.File) error {
	// Skip symbolic link directories
	if f.IsDir() && f.Dir().Ok().IsLink() {
		return g.SkipWalk
	}

	// Skip symbolic link files
	if f.IsLink() {
		return g.SkipWalk
	}

	// Print the path
	f.Path().Ok().Print()

	return nil
}
