package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	f := NewFloat(1.3339)

	f.Hash().MD5().Print()

	String("12.3348992").ToFloat().Unwrap().RoundDecimal(5).Print()

	NewFloat(1.3339).Print()
	NewFloat(13339).Print()

	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))
	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))
}
