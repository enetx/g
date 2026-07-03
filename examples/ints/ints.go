package main

import (
	"fmt"
	"math"

	. "github.com/enetx/g"
)

func main() {
	NewInt(math.MinInt64).RandomRange(math.MaxInt64).Println()

	NewInt(100).Random().Println()
	NewInt(5).RandomRange(-10).Println()

	NewInt(-10).RandomRange(-5).Println()
	NewInt(-10).RandomRange(5).Println()

	NewInt(11).Max(10, 32, 11, 33, 908).Println()
	NewInt(11).Min(-1, 32, 11, 33, 908).Println()

	NewInt(97).Binary().Println()
	NewInt('a').Binary().Println()
	NewInt(byte('a')).Binary().Println()

	// Signum: the sign as -1 / 0 / 1
	fmt.Println(Int(42).Signum(), Int(0).Signum(), Int(-7).Signum()) // 1 0 -1

	// Checked arithmetic — overflow becomes None instead of a silent wraparound
	fmt.Println(Int(math.MaxInt).CheckedAdd(1)) // None
	fmt.Println(Int(10).CheckedDiv(0))          // None — no division panic
	fmt.Println(Int(2).CheckedPow(62))          // Some(4611686018427387904)

	// Saturating / Overflowing variants
	fmt.Println(Int(math.MaxInt).SaturatingAdd(100) == math.MaxInt) // true
	v, overflowed := Int(math.MaxInt).OverflowingAdd(1)
	fmt.Println(v == math.MinInt, overflowed) // true true

	// The checked chain composes with Option combinators
	Int(1_000_000).
		CheckedMul(12).
		Then(func(n Int) Option[Int] { return n.CheckedAdd(500_000) }).
		MapOr("overflow!", func(n Int) string { return n.String().Std() })

}
