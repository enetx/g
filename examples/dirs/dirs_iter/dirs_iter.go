package main

import (
	. "github.com/enetx/g"

	"github.com/enetx/g/fs"
)

func main() {
	// Iterate over and print the names of files in the current directory
	fs.NewDir(".").Read().ForEach(func(f Result[*fs.File]) { f.Ok().Name().Println() })

	// Iterate over and print the full paths of files in the current directory
	fs.NewDir(".").Read().ForEach(func(f Result[*fs.File]) { f.Ok().Path().Ok().Println() })

	// Iterate over and print the names of files in the current directory with a *.go extension
	fs.NewDir("*.go").Glob().ForEach(func(f Result[*fs.File]) { f.Ok().Name().Println() })

	// Iterate over and print the full paths of files in the current directory with a *.go extension
	fs.NewDir("*.go").Glob().ForEach(func(f Result[*fs.File]) { f.Ok().Path().Ok().Println() })

	// Iterate over and print the full paths of non-directory files in the current directory
	fs.NewDir(".").Read().ForEach(func(f Result[*fs.File]) {
		if !f.Ok().Stat().Unwrap().IsDir() {
			f.Ok().Path().Unwrap().Println()
		}
	})
}
