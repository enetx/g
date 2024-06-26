package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	f := g.NewFloat(1.3339)

	f.Hash().MD5().Print()

	g.String("12.3348992").ToFloat().Unwrap().RoundDecimal(5).Print()

	g.NewFloat(1.3339).Print()
	g.NewFloat(13339).Print()

	fmt.Println(g.NewFloat(20).Eq(g.NewFloat(20.0)))
	fmt.Println(g.NewFloat(20).Eq(g.NewFloat(20.0)))
}
