package main

import . "github.com/enetx/g"

func main() {
	// Iterate over and print the names of files in the current directory
	NewDir(".").Read().ForEach(func(f Result[*File]) { f.Ok().Name().Print() })

	// Iterate over and print the full paths of files in the current directory
	NewDir(".").Read().ForEach(func(f Result[*File]) { f.Ok().Path().Ok().Print() })

	// Iterate over and print the names of files in the current directory with a *.go extension
	NewDir("*.go").Glob().ForEach(func(f Result[*File]) { f.Ok().Name().Print() })

	// Iterate over and print the full paths of files in the current directory with a *.go extension
	NewDir("*.go").Glob().ForEach(func(f Result[*File]) { f.Ok().Path().Ok().Print() })

	// Iterate over and print the full paths of non-directory files in the current directory
	NewDir(".").Read().ForEach(func(f Result[*File]) {
		if !f.Ok().Stat().Unwrap().IsDir() {
			f.Ok().Path().Unwrap().Print()
		}
	})
}
