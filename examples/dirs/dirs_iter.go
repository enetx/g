package main

import "github.com/enetx/g"

func main() {
	// Iterate over and print the names of files in the current directory
	g.NewDir(".").Read().Ok().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the full paths of files in the current directory
	g.NewDir(".").Read().Ok().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	// Iterate over and print the names of files in the current directory with a *.go extension
	g.NewDir("*.go").Glob().Unwrap().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the full paths of files in the current directory with a *.go extension
	g.NewDir("*.go").Glob().Unwrap().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	// Iterate over and print the full paths of non-directory files in the current directory
	g.NewDir(".").Read().Unwrap().ForEach(func(f *g.File) {
		if !f.Stat().Unwrap().IsDir() {
			f.Path().Unwrap().Print()
		}
	})
}
