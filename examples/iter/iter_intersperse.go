package main

import (
	"github.com/enetx/g"
)

func main() {
	// Example 1
	g.Slice[string]{"Hello", "World", "!"}.Iter().
		Intersperse(" ").
		Collect().
		Join().
		Print() // Hello World !

		// Example 2
	str := g.String("I love ice cream. Ice cream is delicious.")
	matches := g.Slice[g.String]{"Ice", "cream"}.Iter().Intersperse("").Collect().Append("")

	str = str.ReplaceMulti(matches...)
	str.Print() // I love ice .   is delicious.
}
