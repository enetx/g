package main

import . "github.com/enetx/g"

func main() {
	NewInt(100).Random().Println()

	NewInt(-10).RandomRange(-5).Println()
	NewInt(-10).RandomRange(5).Println()

	NewInt(11).Max(10, 32, 11, 33, 908).Println()
	NewInt(11).Min(-1, 32, 11, 33, 908).Println()

	NewInt(97).Binary().Println()
	NewInt('a').Binary().Println()
	NewInt(byte('a')).Binary().Println()
}
