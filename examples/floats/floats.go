package main

import (
	"fmt"
	"math"

	. "github.com/enetx/g"
)

func main() {
	String("12.3348992").TryFloat().Unwrap().RoundDecimal(5).Println()

	NewFloat(1.3339).Println()
	NewFloat(13339).Println()

	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))
	fmt.Println(NewFloat(20).Eq(NewFloat(20.0)))

	// Float math & classification
	fmt.Println(Float(1).Div(0).IsInf())            // true
	fmt.Println(Float(math.NaN()).Signum().IsNaN()) // true
	fmt.Println(Float(-3.75).Fract())               // -0.75
	Float(2.5).Clamp(0, 2).Println()                // 2
	Float(math.Pi).ToDegrees().Println()            // 180
	Float(0.1).MulAdd(10, -1).Println()             // 5.551115123125783e-17 — true FMA, single rounding

	// Sign-bit checks — signed zero keeps its sign
	fmt.Println(Float(1.5).IsSignPositive())  // true
	fmt.Println(Float(-1.5).IsSignNegative()) // true

	negZero := Float(math.Copysign(0, -1))                     // the Go constant -0.0 is +0, so build -0.0 via Copysign
	fmt.Println(negZero.IsSignNegative(), negZero.Eq(0))       // true true — sign bit set, yet -0.0 == 0.0
	fmt.Println(Float(0.0).IsSignPositive(), negZero.String()) // true -0

}
