package main

import "gitlab.com/x0xO/g"

func main() {
	g.NewDir("./").Read(true).Unwrap().Iter().ForEach(walker)
}

func walker(f *g.File) {
	if f.Stat().Unwrap().IsDir() {
		f.Dir().Unwrap().Read(true).Unwrap().Iter().ForEach(walker)
	}

	f.Path().Unwrap().Print()
}
