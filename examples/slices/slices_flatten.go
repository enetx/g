package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	words := g.SliceOf[g.String]("alpha", "beta", "gamma", "💛💚💙💜", "世界")

	g.MapSlice(words, g.String.Chars).
		Flatten().
		Join().
		Print()
}
