package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	// Iterate over and print the names of files in the current directory
	g.NewDir(".").Read().Ok().Iter().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the names of files in the current directory with full path
	g.NewDir(".").Read(true).Ok().Iter().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the full paths of files in the current directory
	g.NewDir(".").Read().Ok().Iter().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	// Iterate over and print the names of files in the current directory with a *.go extension
	g.NewDir("*.go").Glob().Unwrap().Iter().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the names of files in the current directory with a *.go extension and full path
	g.NewDir("*.go").Glob(true).Unwrap().Iter().ForEach(func(f *g.File) { f.Name().Print() })

	// Iterate over and print the full paths of files in the current directory with a *.go extension
	g.NewDir("*.go").Glob().Unwrap().Iter().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	// Iterate over and print the paths of non-directory files in the current directory with full path
	g.NewDir(".").Read(true).Unwrap().Iter().ForEach(func(f *g.File) {
		if !f.Stat().Unwrap().IsDir() {
			f.Path().Unwrap().Print()
		}
	})
}
