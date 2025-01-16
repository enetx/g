package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	f := NewFloat(1.3339)

	f.Hash().MD5().Println()

	String("12.3348992").ToFloat().Unwrap().RoundDecimal(5).Println()

	NewFloat(1.3339).Println()
	NewFloat(13339).Println()

	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))
	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))
}
