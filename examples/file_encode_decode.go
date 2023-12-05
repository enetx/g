package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	g.NewFile("test.gob").Enc().Gob(g.SliceOf(1, 2, 3)).Unwrap()

	var gobdata g.Slice[int]
	g.NewFile("test.gob").Dec().Gob(&gobdata)

	fmt.Println(gobdata)

	///////////////////////////////////////////////////////////////////

	g.NewFile("test.json").Enc().JSON(g.SliceOf(1, 2, 3)).Unwrap()

	var jsondata g.Slice[int]
	g.NewFile("test.json").Dec().JSON(&jsondata)

	fmt.Println(jsondata)
}
