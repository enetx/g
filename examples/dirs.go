package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	g.NewDir(".").Read().Ok().Iter().ForEach(func(f *g.File) { f.Name().Print() })
	g.NewDir(".").Read(true).Ok().Iter().ForEach(func(f *g.File) { f.Name().Print() })
	g.NewDir(".").Read().Ok().Iter().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	g.NewDir("*.go").Glob().Unwrap().Iter().ForEach(func(f *g.File) { f.Name().Print() })
	g.NewDir("*.go").Glob(true).Unwrap().Iter().ForEach(func(f *g.File) { f.Name().Print() })
	g.NewDir("*.go").Glob().Unwrap().Iter().ForEach(func(f *g.File) { f.Path().Ok().Print() })

	fmt.Println("++++++++++++++")

	g.NewDir(".").Read(true).Unwrap().Iter().ForEach(func(f *g.File) {
		if !f.Stat().Unwrap().IsDir() {
			f.Path().Unwrap().Print()
		}
	})

	d := g.NewDir("")

	d.Path().Unwrap().Print()

	fmt.Println(d.Exist())

	g.NewDir("./some/dir/that/dont/exist/").CreateAll().Unwrap()

	g.NewDir("./some/dir/that/").Rename("./some/dir/aaa/ccc/").Unwrap().Print()

	g.NewDir("aaa").CreateAll().Unwrap().Rename("bbb").Unwrap()

	// make tmp dir
	d = g.NewDir("").CreateTemp().Unwrap().Print()
	d.Remove().Unwrap()

	fmt.Println(d.Exist())
}
