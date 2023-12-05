package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	f := g.NewFile("").CreateTemp("./", "*.txt").Unwrap().Write("some text").Unwrap()
	// f := g.NewFile("").CreateTemp().Unwrap().Write("some text").Unwrap()
	fmt.Println(f.Path().Unwrap(), f.Read().Unwrap())

	f.Read().Unwrap().Hash().MD5().Print()

	fmt.Println(f.Remove().Unwrap().Exist())
}
